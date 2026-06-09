CREATE TABLE `crontab_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `pid` int unsigned NOT NULL DEFAULT '0' COMMENT '进程id',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '名称',
  `exec_cmd` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '运行命令',
  `params` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '参数',
  `result` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '结果',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '状态',
  `start_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '开始时间（毫秒）',
  `end_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '结束时间（毫秒）',
  `cost_time` int unsigned NOT NULL DEFAULT '0' COMMENT '运行时间（毫秒）',
  `partition` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '运行机器',
  `uk` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标识',
  `create_time` timestamp NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_uk_time` (`uk`,`create_time`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='定时任务执行日志';

CREATE TABLE `delay_queue_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '名称',
  `exec_cmd` varchar(50) NOT NULL DEFAULT '' COMMENT '执行方法',
  `params` varchar(500) NOT NULL DEFAULT '' COMMENT '参数',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '执行状态（1待执行 2执行中 3已完成 4失败）',
  `result` varchar(500) NOT NULL DEFAULT '' COMMENT '执行结果',
  `exec_time` int unsigned NOT NULL DEFAULT '0' COMMENT '执行时间（秒）',
  `create_time` timestamp NULL DEFAULT NULL,
  `update_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_status_exectime` (`status`,`exec_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='延时队列记录表';

CREATE TABLE `export_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '标题',
  `header` varchar(2000) NOT NULL COMMENT '表头',
  `query` text NOT NULL COMMENT '请求参数',
  `enums` varchar(2000) NOT NULL DEFAULT '' COMMENT '枚举值',
  `file_name` varchar(255) NOT NULL DEFAULT '' COMMENT '文件名',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态 1 等待导出 2进行中 3已完成 4失败',
  `token` varchar(500) NOT NULL DEFAULT '' COMMENT 'token',
  `create_user_id` int unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
  `file_path` varchar(255) NOT NULL DEFAULT '' COMMENT '文件地址',
  `finish_time` timestamp NULL DEFAULT NULL,
  `error_msg` varchar(255) NOT NULL DEFAULT '' COMMENT '错误信息',
  `create_time` timestamp NULL DEFAULT NULL,
  `update_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='导出日志表';

CREATE TABLE `script_config` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `type` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '类型 1异步消费 2定时任务',
  `cron` varchar(50) NOT NULL DEFAULT '' COMMENT '定时表达式',
  `flag` varchar(50) NOT NULL DEFAULT '' COMMENT 'kafka实例标识',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '名称',
  `exec_cmd` varchar(50) NOT NULL DEFAULT '' COMMENT '执行指令',
  `params` varchar(255) NOT NULL DEFAULT '' COMMENT '参数',
  `topic` varchar(100) NOT NULL DEFAULT '' COMMENT '主题',
  `group_id` varchar(100) NOT NULL DEFAULT '' COMMENT '消费组',
  `key` varchar(100) NOT NULL COMMENT 'redis队列key',
  `progress` tinyint unsigned NOT NULL DEFAULT '2' COMMENT '初始消费者数量',
  `max_progress` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '最大消费者数量',
  `progress_lag_limit` int unsigned NOT NULL DEFAULT '0' COMMENT '消息堆积阈值',
  `progress_avg_msgcount` int unsigned NOT NULL DEFAULT '0' COMMENT '新增单个消费者处理的消息数',
  `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态 1 执行 0不执行',
  `machine` varchar(20) NOT NULL DEFAULT '' COMMENT '运行机器标识',
  `is_log` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '定时任务是否记日志',
  `create_time` timestamp NULL DEFAULT NULL,
  `update_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_machine_staus` (`machine`,`type`,`status`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=71 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='脚本配置表';