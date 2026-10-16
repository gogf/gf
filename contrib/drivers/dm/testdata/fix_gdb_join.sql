DROP TABLE IF EXISTS common_resource;
CREATE TABLE common_resource (
    id BIGINT NOT NULL AUTO_INCREMENT,
    app_id BIGINT NOT NULL,
    resource_id VARCHAR(64) NOT NULL,
    src_instance_id VARCHAR(64) DEFAULT NULL,
    region VARCHAR(36) DEFAULT NULL,
    zone VARCHAR(36) DEFAULT NULL,
    database_kind VARCHAR(20) NOT NULL,
    source_type VARCHAR(64) NOT NULL,
    ip VARCHAR(64) DEFAULT NULL,
    port INT DEFAULT NULL,
    vpc_id VARCHAR(20) DEFAULT NULL,
    subnet_id VARCHAR(20) DEFAULT NULL,
    proxy_ip VARCHAR(64) DEFAULT NULL,
    proxy_port INT DEFAULT NULL,
    proxy_id BIGINT DEFAULT NULL,
    proxy_snat_ip VARCHAR(64) DEFAULT NULL,
    lease_at TIMESTAMP NULL DEFAULT NULL,
    uin VARCHAR(32) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (app_id, src_instance_id, vpc_id, subnet_id, ip, port)
);
INSERT INTO common_resource VALUES (1,1,'2','2','2','3','1','1','1',1,'1','1','1',1,1,'1',NULL,''),(3,2,'3','3','3','3','3','3','3',3,'3','3','3',3,3,'3',NULL,''),(18,1303697168,'dmc-rgnh9qre','vdb-6b6m3u1u','ap-guangzhou','','vdb','cloud','10.0.1.16',80,'vpc-m3dchft7','subnet-9as3a3z2','9.27.72.189',11131,228476,'169.254.128.5, ','2023-11-08 08:13:04',''),(20,1303697168,'dmc-4grzi4jg','tdsqlshard-313spncx','ap-guangzhou','','tdsql','cloud','10.255.0.27',3306,'vpc-407k0e8x','subnet-qhkkk3bo','30.86.239.200',24087,0,'',NULL,'');

DROP TABLE IF EXISTS managed_resource;
CREATE TABLE managed_resource (
    id BIGINT NOT NULL AUTO_INCREMENT,
    instance_id VARCHAR(64) NOT NULL,
    resource_id VARCHAR(64) NOT NULL,
    resource_name VARCHAR(64) DEFAULT NULL,
    status VARCHAR(36) NOT NULL DEFAULT 'valid',
    status_message VARCHAR(64) DEFAULT NULL,
    user_name VARCHAR(64) NOT NULL,
    password VARCHAR(1024) NOT NULL,
    pay_mode TINYINT DEFAULT '0',
    safe_publication BIT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMP NULL DEFAULT NULL,
    deleted TINYINT NOT NULL DEFAULT '0',
    resource_mark_id INT DEFAULT NULL,
    comments VARCHAR(64) DEFAULT NULL,
    rule_template_id VARCHAR(64) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (resource_id)
);
INSERT INTO managed_resource VALUES (1,'2','3','1','1','1','1','1',1,1,'2023-11-06 12:14:21','2023-11-06 12:14:21',NULL,1,1,'1',''),(2,'3','2','1','1','1','1','1',1,0,'2023-11-06 12:15:07','2023-11-06 12:15:07',NULL,1,2,'1',''),(5,'dmcins-jxy0x75m','dmc-rgnh9qre','erichmao-vdb-test','invalid','The Ip field is required','root','2e39af3dd1d447e2b1437b40c62c35995fa22b370c7455ff7815dace3a6e8891ccadcfc893fe1342a4102d742bd7a3e603cd0ac1fcdc072d7c0b5be5836ec87306981b629f9b59aedf0316e9504ab172fa1c95756d5b260114e4feaa0b19223fb61cb268cc4818307ed193dbab830cf556b91cde182686eb70f70ea77f69eff66230dec2ce92bd3352cad31abf47597a5cc6a0d638381dc3bae7aa1b142730790a6d4cefdef1bd460061c966ad5008c2b5fc971b7f4d7dddffa5b1456c45e2917763dd8fffb1fa7fc4783feca95dafc9a9f4edf21b0579f76b0a3154f087e3b9a7fc49af8ff92b12e7b03caa865e72e777dd9d35a11910df0d55ead90e47d5f8',1,1,'2023-11-08 08:13:20','2023-11-09 05:31:07',NULL,0,11,NULL,'12345'),(6,'dmcins-erxms6ya','dmc-4grzi4jg','erichmao-vdb-test','invalid','The Ip field is required','leotaowang','641d846cf75bc7944202251d97dca8335f7f149dd4fd911ca5b87c71ef1dc5d0a66c4e5021ef7ad53136cda2fb2567d34e3dd1a7666e3f64ebf532eb2a55d84952aac86b4211f563f7b9da7dd0f88ec288d6680d3513cea0c1b7ad7babb474717f77ebbc9d63bb458adaf982887da9e63df957ffda572c1c3ed187471b99fdc640b45fed76a6d50dc1090eee79b4d94d056c4d43416133481f55bd040759398680104a84d801e6475dcfe919a00859908296747430b728a00c8d54256ae220235a138e0bbf08fe8b6fc8589971436b55bff966154721a91adbdc9c2b6f50ef5849ed77e5b028116abac51584b8d401cd3a88d18df127006358ed33fc3fa6f480',1,1,'2023-11-08 22:15:17','2023-11-09 05:31:07',NULL,0,11,NULL,'12345');

DROP TABLE IF EXISTS rules_template;
CREATE TABLE rules_template (
    id BIGINT NOT NULL AUTO_INCREMENT,
    app_id BIGINT DEFAULT NULL,
    name VARCHAR(255) NOT NULL,
    database_kind VARCHAR(64) DEFAULT NULL,
    is_default TINYINT NOT NULL DEFAULT '0',
    win_rules VARCHAR(2048) DEFAULT NULL,
    inception_rules VARCHAR(2048) DEFAULT NULL,
    auto_exec_rules VARCHAR(2048) DEFAULT NULL,
    order_check_step VARCHAR(2048) DEFAULT NULL,
    template_id VARCHAR(64) NOT NULL DEFAULT '',
    version INT NOT NULL DEFAULT '1',
    deleted TINYINT NOT NULL DEFAULT '0',
    create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_system TINYINT NOT NULL DEFAULT '0',
    uin VARCHAR(64) DEFAULT NULL,
    subAccountUin VARCHAR(64) DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE (template_id),
    UNIQUE (name, app_id, deleted, uin)
);

DROP TABLE IF EXISTS resource_mark;
CREATE TABLE resource_mark (
    id BIGINT NOT NULL AUTO_INCREMENT,
    app_id BIGINT NOT NULL,
    mark_name VARCHAR(64) NOT NULL,
    color VARCHAR(11) NOT NULL,
    creator VARCHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (app_id, mark_name)
);
INSERT INTO resource_mark VALUES (10,1,'test','red','1','2023-11-06 02:45:46','2023-11-06 02:45:46');
