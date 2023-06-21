-- public.add_soup_data definition

-- Drop table

-- DROP TABLE public.add_soup_data;

CREATE TABLE public.add_soup_data (
	id varchar(255) NOT NULL,
	order_sn varchar(255) NOT NULL, -- 订单号
	table_number varchar(15) NULL, -- 桌号
	shop_id varchar(255) NULL, -- 店铺ID
	bottom_id varchar(100) NULL, -- 锅底ID
	"name" varchar(255) NULL, -- 店铺名称
	order_time timestamp(0) NULL,
	soup_start_time timestamp(0) NULL,
	soup_end_time timestamp(0) NULL,
	feeding_start_time timestamp(0) NULL,
	feeding_end_time timestamp(0) NULL,
	turning_pot_end_time timestamp(0) NULL,
	create_at timestamp(0) NULL,
	bottom_pot varchar(100) NULL,
	CONSTRAINT add_soup_data_pkeys PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.add_soup_data.order_sn IS '订单号';
COMMENT ON COLUMN public.add_soup_data.table_number IS '桌号';
COMMENT ON COLUMN public.add_soup_data.shop_id IS '店铺ID';
COMMENT ON COLUMN public.add_soup_data.bottom_id IS '锅底ID';
COMMENT ON COLUMN public.add_soup_data."name" IS '店铺名称';


-- public.materials definition

-- Drop table

-- DROP TABLE public.materials;

CREATE TABLE public.materials (
	id varchar(50) NOT NULL,
	"name" varchar(50) NULL, -- 物料名称
	dosage int8 NULL, -- 用量
	unit varchar(50) NULL, -- 单位
	water_line int8 NULL, -- 加汤水位标准
	station varchar(50) NULL, -- 工位
	recipe_id varchar(100) NULL, -- 配方ID
	pot_type_id varchar(30) NULL, -- 锅型ID
	resource varchar(20) NULL, -- 物料来源
	CONSTRAINT materials_pkey PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.materials."name" IS '物料名称';
COMMENT ON COLUMN public.materials.dosage IS '用量';
COMMENT ON COLUMN public.materials.unit IS '单位';
COMMENT ON COLUMN public.materials.water_line IS '加汤水位标准';
COMMENT ON COLUMN public.materials.station IS '工位';
COMMENT ON COLUMN public.materials.recipe_id IS '配方ID';
COMMENT ON COLUMN public.materials.pot_type_id IS '锅型ID';
COMMENT ON COLUMN public.materials.resource IS '物料来源';


-- public.original_materials definition

-- Drop table

-- DROP TABLE public.original_materials;

CREATE TABLE public.original_materials (
	id varchar(50) NOT NULL,
	"name" varchar(50) NULL, -- 物料名称
	dosage int8 NULL, -- 用量
	unit varchar(50) NULL, -- 单位
	water_line int8 NULL, -- 加汤水位标准
	station varchar(50) NULL, -- 工位
	resource varchar(20) NULL, -- 来源
	CONSTRAINT original_materials_pkey PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.original_materials."name" IS '物料名称';
COMMENT ON COLUMN public.original_materials.dosage IS '用量';
COMMENT ON COLUMN public.original_materials.unit IS '单位';
COMMENT ON COLUMN public.original_materials.water_line IS '加汤水位标准';
COMMENT ON COLUMN public.original_materials.station IS '工位';
COMMENT ON COLUMN public.original_materials.resource IS '来源';


-- public.original_taste definition

-- Drop table

-- DROP TABLE public.original_taste;

CREATE TABLE public.original_taste (
	id varchar(50) NOT NULL,
	"name" varchar(50) NOT NULL, -- 口味名称
	taste_id varchar(255) NOT NULL, -- 口味ID
	material_id varchar(50) NULL, -- 物料ID
	CONSTRAINT original_taste_pkey PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.original_taste."name" IS '口味名称';
COMMENT ON COLUMN public.original_taste.taste_id IS '口味ID';
COMMENT ON COLUMN public.original_taste.material_id IS '物料ID';


-- public.pot_type definition

-- Drop table

-- DROP TABLE public.pot_type;

CREATE TABLE public.pot_type (
	id varchar(64) NOT NULL, -- 锅型ID
	"name" varchar(255) NULL, -- 锅型名称
	image varchar(255) NULL, -- 图片
	create_at int8 NULL, -- 创建时间
	update_at timestamp(0) NULL DEFAULT CURRENT_TIMESTAMP,
	is_del bool NULL DEFAULT false, -- 是否删除
	soup_standard int8 NULL, -- 加汤水位线标准
	pot_type_id varchar(100) NULL, -- 锅型ID
	CONSTRAINT pot_type_pkey PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.pot_type.id IS '锅型ID';
COMMENT ON COLUMN public.pot_type."name" IS '锅型名称';
COMMENT ON COLUMN public.pot_type.image IS '图片';
COMMENT ON COLUMN public.pot_type.create_at IS '创建时间';
COMMENT ON COLUMN public.pot_type.is_del IS '是否删除';
COMMENT ON COLUMN public.pot_type.soup_standard IS '加汤水位线标准';
COMMENT ON COLUMN public.pot_type.pot_type_id IS '锅型ID';


-- public.recipe definition

-- Drop table

-- DROP TABLE public.recipe;

CREATE TABLE public.recipe (
	id varchar(64) NOT NULL,
	bottom_pot_id varchar(32) NOT NULL, -- 锅底ID
	bottom_pot varchar(64) NOT NULL, -- 锅底名称
	pot_type_id varchar(255) NOT NULL, -- 锅型ID
	pot_type_name varchar(255) NOT NULL, -- 锅型名称
	bottom_properties varchar(255) NOT NULL, -- 锅底属性
	soup_standard int8 NOT NULL, -- 加汤标准
	create_at int8 NOT NULL, -- 创建时间
	update_at timestamp(0) NULL DEFAULT CURRENT_TIMESTAMP, -- 更新时间
	delete_at timestamp(0) NULL, -- 删除时间
	is_del bool NULL DEFAULT false, -- 是否删除
	asset_id varchar(20) NULL, -- 分组ID
	taste_materials varchar(255) NULL, -- 口味物料
	CONSTRAINT recipe_pkey PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.recipe.bottom_pot_id IS '锅底ID';
COMMENT ON COLUMN public.recipe.bottom_pot IS '锅底名称';
COMMENT ON COLUMN public.recipe.pot_type_id IS '锅型ID';
COMMENT ON COLUMN public.recipe.pot_type_name IS '锅型名称';
COMMENT ON COLUMN public.recipe.bottom_properties IS '锅底属性';
COMMENT ON COLUMN public.recipe.soup_standard IS '加汤标准';
COMMENT ON COLUMN public.recipe.create_at IS '创建时间';
COMMENT ON COLUMN public.recipe.update_at IS '更新时间';
COMMENT ON COLUMN public.recipe.delete_at IS '删除时间';
COMMENT ON COLUMN public.recipe.is_del IS '是否删除';
COMMENT ON COLUMN public.recipe.asset_id IS '分组ID';
COMMENT ON COLUMN public.recipe.taste_materials IS '口味物料';


-- public.taste definition

-- Drop table

-- DROP TABLE public.taste;

CREATE TABLE public.taste (
	id varchar(50) NOT NULL,
	"name" varchar(50) NOT NULL, -- 口味名称
	taste_id varchar(255) NOT NULL, -- 口味ID
	create_at int8 NOT NULL, -- 创建时间
	update_at timestamp(0) NULL DEFAULT CURRENT_TIMESTAMP, -- 更新时间
	delete_at timestamp(0) NULL, -- 删除时间
	is_del bool NULL, -- 是否删除
	recipe_id varchar(100) NULL, -- 配方ID
	pot_type_id varchar(20) NULL, -- 锅型ID
	material_id_list varchar(2000) NULL, -- 物料 ID
	bottom_pot_id varchar(100) NULL, -- 订单配方id
	CONSTRAINT taste_pkey PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.taste."name" IS '口味名称';
COMMENT ON COLUMN public.taste.taste_id IS '口味ID';
COMMENT ON COLUMN public.taste.create_at IS '创建时间';
COMMENT ON COLUMN public.taste.update_at IS '更新时间';
COMMENT ON COLUMN public.taste.delete_at IS '删除时间';
COMMENT ON COLUMN public.taste.is_del IS '是否删除';
COMMENT ON COLUMN public.taste.recipe_id IS '配方ID';
COMMENT ON COLUMN public.taste.pot_type_id IS '锅型ID';
COMMENT ON COLUMN public.taste.material_id_list IS '物料 ID';
COMMENT ON COLUMN public.taste.bottom_pot_id IS '订单配方id';

INSERT INTO public.tp_function (id, function_name, menu_id, "path", "name", component, title, icon, "type", function_code, parent_id, sort, tenant_id, sys_flag) VALUES('3be98efe-706b-bf27-2c84-bff4cb9d2661', '', NULL, '/recipe/index', 'RecipeList', '/pages/recipe/index/index.vue', '配方管理', '', '1', '', '0', 100, NULL, NULL);
INSERT INTO public.tp_function (id, function_name, menu_id, "path", "name", component, title, icon, "type", function_code, parent_id, sort, tenant_id, sys_flag) VALUES('1dadabe8-659f-4a6c-a9fb-4314579cb3ab', '', NULL, '/pot/index', 'PotIndex', '/pages/pot/index/index.vue', '锅型管理', '', '1', '', '0', 99, NULL, NULL);
INSERT INTO public.tp_function (id, function_name, menu_id, "path", "name", component, title, icon, "type", function_code, parent_id, sort, tenant_id, sys_flag) VALUES('a27e0e68-5263-51a2-5cc4-5ea6130e1aef', '', NULL, '/soup/index', 'SoupDataManage', '/pages/soup/index/index.vue', '加汤数据管理', '', '1', '', '0', 0, NULL, NULL);
