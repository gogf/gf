package main

import (
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/os/glog"
)

func main() {
	type Detail struct {
		Mall_coupon_id                  int64    `json:"mall_coupon_id"`                  //店铺优惠券id
		Mall_coupon_discount_pct        int      `json:"mall_coupon_discount_pct"`        //店铺折扣
		Mall_coupon_min_order_amount    int      `json:"mall_coupon_min_order_amount"`    //最小使用金额
		Mall_coupon_max_discount_amount int      `json:"mall_coupon_max_discount_amount"` //最大使用金额
		Mall_coupon_total_quantity      int64    `json:"mall_coupon_total_quantity"`      //店铺券总量
		Mall_coupon_remain_quantity     int64    `json:"mall_coupon_remain_quantity"`     //店铺券余量
		Mall_coupon_start_time          int64    `json:"mall_coupon_start_time"`          //店铺券使用开始时间
		Mall_coupon_end_time            int64    `json:"mall_coupon_end_time"`            //店铺券使用结束时间
		Goods_id                        int64    `json:"goods_id"`                        //参与多多进宝的商品ID
		Goods_name                      string   `json:"goods_name"`                      //参与多多进宝的商品标题
		Goods_desc                      string   `json:"goods_desc"`                      //参与多多进宝的商品描述
		Goods_image_url                 string   `json:"goods_image_url"`                 //多多进宝商品主图
		Goods_gallery_urls              []string `json:"goods_gallery_urls"`              //商品轮播图
		Min_group_price                 int64    `json:"min_group_price"`                 //最低价sku的拼团价，单位为分
		Min_normal_price                int64    `json:"min_normal_price"`                //最低价sku的单买价，单位为分
		Mall_name                       string   `json:"mall_name"`                       //店铺名称
		Opt_id                          int64    `json:"opt_id"`                          //商品标签ID，使用pdd.goods.opt.get接口获取
		Opt_name                        string   `json:"opt_name"`                        //商品标签名称
		Opt_ids                         []int    `json:"opt_ids"`                         //商品标签ID
		Cat_ids                         []int    `json:"cat_ids"`                         //商品一~四级类目ID列表
		Coupon_min_order_amount         int64    `json:"coupon_min_order_amount"`         //优惠券门槛金额，单位为分
		Coupon_discount                 int64    `json:"coupon_discount"`                 //优惠券面额，单位为分
		Coupon_total_quantity           int64    `json:"coupon_total_quantity"`           //优惠券总数量
		Coupon_remain_quantity          int64    `json:"coupon_remain_quantity"`          //优惠券剩余数量
		Coupon_start_time               int64    `json:"coupon_start_time"`               //优惠券生效时间，UNIX时间戳
		Coupon_end_time                 int64    `json:"coupon_end_time"`                 //优惠券失效时间，UNIX时间戳
		Promotion_rate                  int64    `json:"promotion_rate"`                  //佣金比例，千分比
		Goods_eval_count                int64    `json:"goods_eval_count"`                //商品评价数
		Cat_id                          int64    `json:"cat_id"`                          //商品类目ID，使用pdd.goods.cats.get接口获取
		Sales_tip                       string   `json:"sales_tip"`                       //已售卖件数
		Mall_id                         int64    `json:"mall_id"`                         //商家id
		Service_tags                    []int64  `json:"service_tags"`                    // 服务标签: 4-送货入户并安装,5-送货入户,6-电子发票,9-坏果包赔 ...
		Clt_cpn_batch_sn                string   `json:"clt_cpn_batch_sn"`                //店铺收藏券id
		Clt_cpn_start_time              int64    `json:"clt_cpn_start_time"`              //店铺收藏券起始时间
		Clt_cpn_end_time                int64    `json:"clt_cpn_end_time"`                //店铺收藏券截止时间
		Clt_cpn_quantity                int64    `json:"clt_cpn_quantity"`                //店铺收藏券总量
		Clt_cpn_remain_quantity         int64    `json:"clt_cpn_remain_quantity"`         //店铺收藏券剩余量
		Clt_cpn_discount                int64    `json:"clt_cpn_discount"`                //店铺收藏券面额，单位为分
		Clt_cpn_min_amt                 int64    `json:"clt_cpn_min_amt"`                 //店铺收藏券使用门槛价格，单位为分
		Desc_txt                        string   `json:"desc_txt"`                        //描述分
		Serv_txt                        string   `json:"serv_txt"`                        //服务分
		Lgst_txt                        string   `json:"lgst_txt"`                        //物流分
		Plan_type                       int      `json:"plan_type"`                       //推广计划类型
		Zs_duo_id                       int64    `json:"zs_duo_id"`                       //招商团长id
		Only_scene_auth                 int      `json:"only_scene_auth"`                 //快手专享
	}

	s := `{"goods_detail_response":{"goods_details":[{"category_name":"百货","clt_cpn_end_time":null,"clt_cpn_min_amt":null,"coupon_remain_quantity":7000,"clt_cpn_remain_quantity":null,"promotion_rate":110,"coupon_id":3248825612,"service_tags":[24,13],"mall_id":983152922,"mall_name":"唯美之恋麻雀专卖店","mall_coupon_end_time":0,"clt_cpn_batch_sn":null,"lgst_txt":"高","goods_name":"【40卷32卷12卷】家用实惠卫生纸无芯卷纸批发纸巾厕纸卷筒纸手纸","clt_cpn_discount":null,"goods_id":9133961778,"goods_gallery_urls":["https://t00img.yangkeduo.com/goods/images/2019-08-18/e05846e8-f317-47bf-987d-5ccc0864e60d.jpg","https://t00img.yangkeduo.com/goods/images/2019-12-06/91756ca8-f001-4394-87d3-e7a67dc81d5a.jpg","https://t00img.yangkeduo.com/goods/images/2019-12-07/14c8abe5-cc27-4265-a683-7e2b089bb82e.jpg","https://t00img.yangkeduo.com/goods/images/2019-08-18/c4d338c0-1922-463f-bfcd-accd02018e44.jpg","https://t00img.yangkeduo.com/goods/images/2019-08-18/2b7f67f8-7738-4462-8f27-243594401ecd.jpg","https://t00img.yangkeduo.com/goods/images/2019-08-18/9b9bf86c-a3a7-4de1-9169-8b25c756ed01.jpg","https://t00img.yangkeduo.com/goods/images/2019-12-01/e373de68-c8a4-4752-adbc-d26b1040f8b7.jpg","https://t00img.yangkeduo.com/goods/images/2019-11-25/e622ef70-d917-4606-86af-a238db8c2fce.jpg","https://t00img.yangkeduo.com/goods/images/2019-11-25/9dd92cdd-1a07-4261-9584-cfafb0ddd421.jpg","https://t00img.yangkeduo.com/goods/images/2019-12-03/f7088db8-02f9-412b-a96c-5b452f41d9e3.jpg"],"goods_desc":"本产品原材料精选原生木浆制造,先进喷浆制造工艺,纸质更柔软细腻,没有纸屑,更有韧劲。嘉禾纸业建厂15年,两大自主品牌纸语和唯美之恋,年产能两万万余吨,拥有先进的造纸设备,采用巴西进口制浆制作制造而成,纸质柔软细腻,吸水性强,致力于打造国内孕婴卫生纸首选品牌,深受广大消费者好评。","opt_name":"百货","opt_ids":[17665,10116,10696,330,12619,12,11212,8590,10702,8591,15,13903,8592,17680,12691,22102,13911,12696,22105,21915,223,292,12581,13926,360,10730,12586,20527,11187,8569,12729,17657,8570,8571,10110,10111],"goods_image_url":"https://t00img.yangkeduo.com/goods/images/2019-08-18/e05846e8-f317-47bf-987d-5ccc0864e60d.jpg","has_mall_coupon":false,"min_group_price":720,"coupon_start_time":1575820800,"coupon_discount":200,"coupon_end_time":1575993599,"zs_duo_id":0,"mall_coupon_remain_quantity":0,"plan_type":2,"clt_cpn_quantity":null,"crt_rf_ordr_rto1m":0.03807085638781072,"cat_ids":[17285,17297,17402],"coupon_min_order_amount":200,"category_id":17657,"mall_coupon_discount_pct":0,"cat_id":null,"coupon_total_quantity":12000,"mall_coupon_min_order_amount":0,"merchant_type":4,"clt_cpn_start_time":null,"sales_tip":"4.2万","plan_type_all":4,"only_scene_auth":true,"mall_coupon_id":0,"desc_txt":"高","goods_thumbnail_url":"https://t00img.yangkeduo.com/goods/images/2019-11-20/0b0d078822c6e777ed1eb4ac2cdc930a.jpeg","opt_id":17657,"search_id":null,"min_normal_price":1190,"has_coupon":true,"mall_coupon_start_time":0,"serv_txt":"高","mall_rate":110,"mall_coupon_total_quantity":0,"create_at":null,"mall_coupon_max_discount_amount":0,"mall_cps":1}],"request_id":"15758638812810680"}}`
	if j, err := gjson.DecodeToJson([]byte(s)); err != nil {
		glog.Error(err)
	} else {

		inter := j.GetInterfaces("goods_detail_response.goods_details.0")
		fmt.Println("接口:", inter)
		detail := Detail{}
		err := j.GetStruct("goods_detail_response.goods_details.0", &detail)
		fmt.Println("结构体:", err, detail)
	}
}
