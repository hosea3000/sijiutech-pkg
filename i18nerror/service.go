package i18nerr

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"code.sijiutech.com/ideastudio/sijiutech-pkg/i18nerror/helper"
	"code.sijiutech.com/ideastudio/sijiutech-pkg/i18nerror/sqlc/repository"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

type ServiceName string

type I18NService interface {
	Translate(ctx context.Context, language, key string) (string, error)
}

func NewI18NService(
	dbx repository.DBTX,
	redisClient *redis.Client,
	serviceName ServiceName,
) I18NService {
	repo := repository.New(dbx)
	service := &i18NService{
		repository:  repo,
		redisClient: redisClient,
		serviceName: serviceName,
	}
	SetI18nService(service)
	return service
}

type i18NService struct {
	redisClient *redis.Client
	repository  *repository.Queries
	serviceName ServiceName
}

func (s *i18NService) Translate(ctx context.Context, language, key string) (string, error) {
	translatedValue := key
	var err error

	// 1. 查找不存在的国际化key
	if s.ExistsInNotFundRedis(ctx, language, key) {
		return "", errors.New("i18n not found:" + key)
	}
	// 2. 查找 redis key
	redisKey := fmt.Sprintf("i18n:%s:%s", string(s.serviceName), language)

	// 1. 先从 Redis Hash 中获取特定 key 的翻译值
	translatedValue, err = s.redisClient.HGet(ctx, redisKey, key).Result()
	// translatedValue 如果为 0 说明是富文本内容，富文本不直接返回，去数据库查询
	if err == nil && translatedValue != "" && translatedValue != "0" {
		// 找到了，直接返回
		return translatedValue, nil
	}
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 如果是 key 不存在，继续查询数据库
		} else {
			return "", err
		}
	}

	// 2. 从数据库查询
	dataRow, err := s.TranslatedValue(ctx, language, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err := s.SetInNotFundRedis(ctx, language, key); err != nil {
				return "", err
			}
			return "", errors.New("i18n not found:" + key)
		}
		return "", err
	}
	translatedValue = dataRow.Value.String

	// 如果数据类型为 0，则表示是普通文本，普通文本写进redis，富文本不写
	if dataRow.DataType == 0 {
		// 3. 将查询到的值存储到 Redis 中
		err = s.redisClient.HSet(ctx, redisKey, key, translatedValue).Err()
		if err != nil {
			return "", err
		}
	}

	return translatedValue, nil
}

func (s *i18NService) SetInNotFundRedis(ctx context.Context, lang, key string) error {
	redisNotFoundKey := fmt.Sprintf("i18n:i18n-no-found-message:%s_%s", lang, key)
	return s.redisClient.Set(ctx, redisNotFoundKey, "\""+key+"\"", 5*time.Minute).Err()
}

func (s *i18NService) ExistsInNotFundRedis(ctx context.Context, lang, key string) bool {
	redisNotFoundKey := fmt.Sprintf("i18n:i18n-no-found-message:%s_%s", lang, key)
	return s.redisClient.Exists(ctx, redisNotFoundKey).Val() > 0
}

// getTerminalCodeList 获取终端code集合`
// 参考 Java: getTerminalCodeList(String terminalCode)
// 返回父终端以及所有子终端的 code 列表
func (s *i18NService) getTerminalCodeList(ctx context.Context, terminalCode string) ([]string, error) {
	codeList := []string{terminalCode}

	// 根据 terminalCode 查询终端
	terminal, err := s.repository.SelectLanguageTerminalByCode(ctx, helper.ToSqlNullString(terminalCode))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 如果找不到终端，只返回传入的 code
			return codeList, nil
		}
		return nil, errors.Wrap(err, "select language terminal by code failed")
	}

	// 查询所有子终端（parent_id 等于该终端的 id）
	childTerminals, err := s.repository.SelectLanguageTerminalByParentId(ctx, helper.ToSqlNullInt64(terminal.ID))
	if err != nil {
		return nil, errors.Wrap(err, "select language terminal by parent id failed")
	}

	// 将子终端的 code 添加到列表中
	if len(childTerminals) > 0 {
		for _, child := range childTerminals {
			if child.Code.Valid {
				codeList = append(codeList, child.Code.String)
			}
		}
	}

	return codeList, nil
}

func (s *i18NService) TranslatedValue(ctx context.Context, lang, key string) (*repository.SelectI18nByLangAndKeyRow, error) {
	// 先查询出所有的子终端
	terminalCodeList, err := s.getTerminalCodeList(ctx, string(s.serviceName))
	if err != nil {
		return nil, errors.Wrap(err, "get terminal code list failed")
	}

	// 查询翻译（需要遍历所有终端代码）
	terminalCodes := lo.Map(terminalCodeList, func(item string, index int) sql.NullString {
		return helper.ToSqlNullString(item)
	})
	dataRow, err := s.repository.SelectI18nByLangAndKey(ctx, repository.SelectI18nByLangAndKeyParams{
		Lang:          helper.ToSqlNullString(lang),
		Key:           key,
		TerminalCodes: terminalCodes,
	})
	if err != nil {
		return nil, errors.Wrap(err, "select i18n by lang and key failed")
	}

	return &dataRow, nil
}
