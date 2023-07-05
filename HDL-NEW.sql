CREATE TABLE public.hdl_materials (
	id varchar(50) NOT NULL,
	"name" varchar(500) NULL, -- 物料名称
	dosage int8 NULL, -- 用量
	unit varchar(50) NULL, -- 单位
	water_line int8 NULL, -- 加汤水位标准
	station varchar(50) NULL, -- 工位（鲜料工位、传锅工位、所有工位）
	resource varchar(50) NULL, -- 物料来源material-锅底 taste-口味
	remark varchar(500) NULL, -- 保留
	CONSTRAINT hdl_materials_pkey PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.hdl_materials."name" IS '物料名称';
COMMENT ON COLUMN public.hdl_materials.dosage IS '用量';
COMMENT ON COLUMN public.hdl_materials.unit IS '单位';
COMMENT ON COLUMN public.hdl_materials.water_line IS '加汤水位标准';
COMMENT ON COLUMN public.hdl_materials.station IS '工位';
COMMENT ON COLUMN public.hdl_materials.resource IS '物料来源';


CREATE TABLE public.hdl_pot_type (
	id varchar(36) NOT NULL, -- 锅型ID
	"name" varchar(500) NULL, -- 锅型名称
	image varchar(500) NULL, -- 图片
	create_at int8 NULL, -- 创建时间
	update_at int8 NULL,
	soup_standard int8 NULL, -- 加汤水位线标准
	pot_type_id varchar(50) NULL, -- 锅型ID
	remark varchar(500) NULL, -- 锅型ID
	CONSTRAINT pot_type_pkey PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.hdl_pot_type.id IS '锅型ID';
COMMENT ON COLUMN public.hdl_pot_type."name" IS '锅型名称';
COMMENT ON COLUMN public.hdl_pot_type.image IS '图片';
COMMENT ON COLUMN public.hdl_pot_type.create_at IS '创建时间';
COMMENT ON COLUMN public.hdl_pot_type.soup_standard IS '加汤水位线标准';
COMMENT ON COLUMN public.hdl_pot_type.pot_type_id IS '锅型ID';

CREATE TABLE public.hdl_taste (
	id varchar(36) NOT NULL,
	"name" varchar(500) NOT NULL, -- 口味名称
	taste_id varchar(50) NOT NULL, -- 口味ID
	create_at int8 NULL, -- 
	update_time int8 NULL, -- 
	remark  varchar(500) NOT NULL, 
	CONSTRAINT hdl_taste_pkey PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.hdl_taste."name" IS '口味名称';
COMMENT ON COLUMN public.hdl_taste.taste_id IS '口味ID';

CREATE TABLE public.hdl_recipe (
	id varchar(36) NOT NULL,
	bottom_pot_id varchar(500) NOT NULL, -- 锅底ID
	bottom_pot varchar(500) NOT NULL, -- 锅底名称
	bottom_properties varchar(50) NOT NULL, -- 锅底属性
	hdl_pot_type_id varchar(36) NOT NULL, -- 锅型ID
	create_at int8 NOT NULL, -- 创建时间
	update_at int8 NULL , -- 更新时间
	tenant_id  varchar(36) NOT NULL, -- 租户id
	remark varchar(500) NULL, -- 分组ID
	CONSTRAINT hdl_recipe_pkey PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.hdl_recipe.bottom_pot_id IS '锅底ID';
COMMENT ON COLUMN public.hdl_recipe.bottom_pot IS '锅底名称';
COMMENT ON COLUMN public.hdl_recipe.hdl_pot_type_id IS '锅型ID';
COMMENT ON COLUMN public.hdl_recipe.bottom_properties IS '锅底属性';
COMMENT ON COLUMN public.hdl_recipe.create_at IS '创建时间';
COMMENT ON COLUMN public.hdl_recipe.update_at IS '更新时间';


CREATE TABLE public.hdl_r_recipe_materals  (
	hdl_recipe_id varchar(36) NOT NULL,
	hdl_materals_id varchar(36) NOT NULL
);

CREATE TABLE public.hdl_r_taste_materals  (
	hdl_taste_id varchar(36) NOT NULL,
	hdl_materals_id varchar(36) NOT NULL
);

CREATE TABLE public.hdl_r_recipe_taste   (
	hdl_taste_id varchar(36) NOT NULL,
	hdl_recipe_id varchar(36) NOT NULL
);



CREATE TABLE public.hdl_add_soup_data (
	id varchar(36) NOT NULL,
	order_sn varchar(200) NOT NULL, -- 订单号
	table_number varchar(200) NULL, -- 桌号
	shop_id varchar(50) NULL, -- 店铺ID
	bottom_id varchar(200) NULL, -- 锅底ID
	"name" varchar(200) NULL, -- 店铺名称
	order_time timestamp(100) NULL, -- 订单时间
	soup_start_time timestamp(100) NULL,-- 加汤开始时间
	soup_end_time timestamp(100) NULL, -- 加汤结束时间
	feeding_start_time timestamp(100) NULL, --加料开始时间
	feeding_end_time timestamp(100) NULL, -- 加料结束时间
	turning_pot_end_time timestamp(100) NULL,-- 转锅结束时间
	create_at int8 NULL,
	bottom_pot varchar(100) NULL,
	CONSTRAINT hdl_add_soup_data_pkeys PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.hdl_add_soup_data.order_sn IS '订单号';
COMMENT ON COLUMN public.hdl_add_soup_data.table_number IS '桌号';
COMMENT ON COLUMN public.hdl_add_soup_data.shop_id IS '店铺ID';
COMMENT ON COLUMN public.hdl_add_soup_data.bottom_id IS '锅底ID';
COMMENT ON COLUMN public.hdl_add_soup_data."name" IS '店铺名称';
COMMENT ON COLUMN public.hdl_add_soup_data.order_time IS '订单时间';
COMMENT ON COLUMN public.hdl_add_soup_data.soup_start_time IS '加汤开始时间';
COMMENT ON COLUMN public.hdl_add_soup_data.soup_end_time IS '加汤结束时间';
COMMENT ON COLUMN public.hdl_add_soup_data.feeding_start_time IS '加料开始时间';
COMMENT ON COLUMN public.hdl_add_soup_data.feeding_end_time IS '加料结束时间';
COMMENT ON COLUMN public.hdl_add_soup_data.turning_pot_end_time IS '转锅结束时间';

ALTER TABLE public.hdl_materials ADD create_at int8 NULL;
ALTER TABLE public.hdl_r_recipe_materals RENAME TO hdl_r_recipe_materials;
ALTER TABLE public.hdl_r_taste_materals RENAME TO hdl_r_taste_materials;
ALTER TABLE public.hdl_r_recipe_materials RENAME COLUMN hdl_materals_id TO hdl_materials_id;
ALTER TABLE public.hdl_r_taste_materials RENAME COLUMN hdl_materals_id TO hdl_materials_id;
ALTER TABLE public.hdl_r_recipe_materials ADD CONSTRAINT hdl_r_recipe_materials_fk FOREIGN KEY (hdl_recipe_id) REFERENCES public.hdl_recipe(id) ON DELETE CASCADE;
ALTER TABLE public.hdl_r_recipe_materials ADD CONSTRAINT hdl_r_recipe_materials_fk_1 FOREIGN KEY (hdl_materials_id) REFERENCES public.hdl_materials(id) ON DELETE RESTRICT;
ALTER TABLE public.hdl_r_recipe_taste ADD CONSTRAINT hdl_r_recipe_taste_fk FOREIGN KEY (hdl_recipe_id) REFERENCES public.hdl_recipe(id) ON DELETE CASCADE;
ALTER TABLE public.hdl_r_recipe_taste ADD CONSTRAINT hdl_r_recipe_taste_fk_1 FOREIGN KEY (hdl_taste_id) REFERENCES public.hdl_taste(id) ON DELETE RESTRICT;
ALTER TABLE public.hdl_r_taste_materials ADD CONSTRAINT hdl_r_taste_materials_fk FOREIGN KEY (hdl_taste_id) REFERENCES public.hdl_taste(id) ON DELETE CASCADE;
ALTER TABLE public.hdl_r_taste_materials ADD CONSTRAINT hdl_r_taste_materials_fk_1 FOREIGN KEY (hdl_materials_id) REFERENCES public.hdl_materials(id) ON DELETE RESTRICT;

INSERT INTO public.tp_function (id, function_name, menu_id, "path", "name", component, title, icon, "type", function_code, parent_id, sort, tenant_id, sys_flag) VALUES('3be98efe-706b-bf27-2c84-bff4cb9d2661', '', NULL, '/recipe/index', 'RecipeList', '/pages/recipe/index.vue', '配方管理', '', '1', '', '0', 100, NULL, NULL);
INSERT INTO public.tp_function (id, function_name, menu_id, "path", "name", component, title, icon, "type", function_code, parent_id, sort, tenant_id, sys_flag) VALUES('1dadabe8-659f-4a6c-a9fb-4314579cb3ab', '', NULL, '/pot/index', 'PotIndex', '/pages/pot/index.vue', '锅型管理', '', '1', '', '0', 99, NULL, NULL);
INSERT INTO public.tp_function (id, function_name, menu_id, "path", "name", component, title, icon, "type", function_code, parent_id, sort, tenant_id, sys_flag) VALUES('a27e0e68-5263-51a2-5cc4-5ea6130e1aef', '', NULL, '/soup/index', 'SoupDataManage', '/pages/soup/index.vue', '加汤数据管理', '', '1', '', '0', 0, NULL, NULL);
