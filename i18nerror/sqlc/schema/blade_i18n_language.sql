CREATE TABLE `blade_i18n_language`
(
    `id`          bigint  NOT NULL COMMENT 'id',
    `code`        varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '编码',
    `name`        varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '名称',
    `is_default`  tinyint NOT NULL                                             DEFAULT '0' COMMENT '0-否 1-是   是否系统默认语种',
    `create_user` bigint                                                       DEFAULT NULL COMMENT '创建人',
    `create_time` datetime                                                     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_user` bigint                                                       DEFAULT NULL COMMENT '修改人',
    `update_time` datetime                                                     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `is_deleted`  int                                                          DEFAULT '0' COMMENT '是否已删除',
    `status`      int                                                          DEFAULT NULL COMMENT '状态',
    `create_dept` bigint                                                       DEFAULT NULL COMMENT '创建部门',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='i18n语言定义'

