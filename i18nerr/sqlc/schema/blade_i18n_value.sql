CREATE TABLE `blade_i18n_value`
(
    `id`            bigint NOT NULL COMMENT 'id',
    `definition_id` bigint   DEFAULT NULL COMMENT '语言定义id',
    `key_id`        bigint   DEFAULT NULL COMMENT '语言键id',
    `value`         text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '键值',
    `create_user`   bigint   DEFAULT NULL COMMENT '创建人',
    `create_time`   datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_user`   bigint   DEFAULT NULL COMMENT '修改人',
    `update_time`   datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `create_dept`   bigint   DEFAULT NULL COMMENT '创建部门',
    `status`        int      DEFAULT NULL COMMENT '状态',
    `is_deleted`    int      DEFAULT '0' COMMENT '是否已删除',
    PRIMARY KEY (`id`) USING BTREE,
    KEY             `idx_definition_id` (`definition_id`) USING BTREE,
    KEY             `idx_key_id` (`key_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='i18n语言值'

