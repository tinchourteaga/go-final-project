drop database if exists melisprint;
create database melisprint;
use melisprint;
create table localities(
    `id` varchar(10) not null primary key,
    locality_name text not null,
    province_name text not null,
    country_name text not null
);
create table sellers(
    `id` int not null primary key auto_increment,
    cid int not null,
    company_name text not null,
    `address` text not null,
    telephone varchar(15) not null,
    locality_id varchar(10) not null,
    foreign key (locality_id) references localities(id)
);
create table products(
    `id` int not null primary key auto_increment,
    `description` text not null,
    expiration_rate int not null,
    freezing_rate int not null,
    height float not null,
    lenght float not null,
    netweight float not null,
    product_code text not null,
    recommended_freezing_temperature float not null,
    width float not null,
    id_product_type int not null,
    id_seller int,
    foreign key (id_seller) references sellers(id)
);
create table warehouses(
    `id` int not null primary key auto_increment,
    `address` text null,
    telephone text null,
    warehouse_code text null,
    minimum_capacity int null,
    minimum_temperature int null
);
create table employees(
    `id` int not null primary key auto_increment,
    card_number_id text not null,
    first_name text not null,
    last_name text not null,
    warehouse_id int not null,
    foreign key (warehouse_id) references warehouses(id)
);
create table sections(
    `id` int not null primary key auto_increment,
    section_number int not null,
    current_temperature int not null,
    minimum_temperature int not null,
    current_capacity int not null,
    minimum_capacity int not null,
    maximum_capacity int not null,
    warehouse_id int not null,
    id_product_type int not null,
    foreign key (warehouse_id) references warehouses(id)
);
create table buyers(
    `id` int not null primary key auto_increment,
    card_number_id text not null,
    first_name text not null,
    last_name text not null
);
create table product_records(
	`id` int not null primary key auto_increment,
    last_update_date date,
    purchase_price float,
    sale_price float,
    product_id int not null,
    foreign key (product_id) references products(id)
);
create table product_batches(
    id int not null primary key auto_increment,
    batch_number int not null unique,
    current_quantity int not null,
    current_temperature int not null,
    due_date date not null,
    initial_quantity int not null,
    manufacturing_date date not null,
    manufacturing_hour int not null,
	minimum_temperature int not null,
    product_id int not null ,
    foreign key (product_id) references products(id),
    section_id int not null,
    foreign key (section_id) references sections(id)
);
create table carries(
    `id` int not null primary key auto_increment,
    cid varchar(10) unique null,
    company_name text null,
    `address` text null,
    telephone text null,
    locality_id varchar(10) null,
    foreign key (locality_id) references localities(id)
);
create table inbound_orders(
    `id` int not null primary key auto_increment,
    order_date text not null,
    order_number varchar(100) not null unique,
    employee_id int not null,
    product_batch_id int not null,
    warehouse_id int not null,
    foreign key (employee_id) references employees(id),
    foreign key (warehouse_id) references warehouses(id),
    foreign key (product_batch_id) references product_batches(id)
);
create table purchase_orders(
	`id` int not null primary key auto_increment,
    order_number text not null,
    order_date date not null,
    tracking_code text not null,
    buyer_id int not null,
    product_record_id int not null,
    order_status_id int not null,
    foreign key (product_record_id) references product_records(id),
	foreign key (buyer_id) references buyers(id)
);
create table logs(
    `id` int not null primary key auto_increment,
    time_stamp text not null,
    `user` text not null,
    file_path text not null,
    function_line text not null,
    caller_function text not null,
    msg text not null
);
