package whole

import(
    // "github.com/fecshopsoft/fec-go/shell/whole"
    // "github.com/fecshopsoft/fec-go/util"
    "github.com/fecshopsoft/fec-go/db/mongodb"
    "github.com/fecshopsoft/fec-go/db/esdb"
    "github.com/fecshopsoft/fec-go/helper"
    "github.com/fecshopsoft/fec-go/shell/model"
    "github.com/globalsign/mgo"
    "github.com/globalsign/mgo/bson"
    "math"
    "log"
    
)

// 全部统计计算部分
func AllMapReduct(dbName string, collName string, outCollName string, website_id string) error {
    var err error
    mapStr := `
        function() {  
            // browser_name = this.browser_name ? this.browser_name : null;
            // login_email 	= this.login_email ? [this.login_email] : null;
            website_id = "` + website_id + `"
            stay_seconds = this.stay_seconds ? this.stay_seconds : 0;
            
            // 跳出个数和退出个数
            jump_out_count = 0;   //跳出
            drop_out_count = 0;		//退出
            // 最后一个页面
            rate_pv = 0;
            if(this.stay_seconds == 0){
                // 如果这个页面是第一个访问的页面，则代表跳出
                if(this.uuid_first_page == 1 ){
                    jump_out_count = 1;
                }else{
                    drop_out_count = 1;
                }
            }else{
                rate_pv = 1;
            }
            
            // browser_name
            browser_name  = this.browser_name;
            if(browser_name){
                b= {};
                b[browser_name] = 1;
                browser_name 	= b;
            }else{
                browser_name = null;
            }
            
            // devide
            devide 			= this.devide;
            if(devide){
                b= {};
                b[devide] = 1;
                devide 	= b;
            }else{
                devide = null;
            }
            
            // country_code
            country_code 	= this.country_code;
            if(country_code){
                b= {};
                b[country_code] = 1;
                country_code 	= b;
            }else{
                country_code = null;
            }
            
            // operate
            operate 		= this.operate ;
            if(operate){
                b= {};
                b[operate] = 1;
                operate 	= b;
            }else{
                operate = null;
            }
            // fec_app
            fec_app 		= this.fec_app ;
            if(fec_app){
                b= {};
                b[fec_app] = 1;
                fec_app 	= b;
            }else{
                fec_app = null;
            }
            
            resolution 		= this.resolution ;
            if(resolution){
                b= {};
                b[resolution] = 1;
                resolution 	= b;
            }else{
                resolution = null;
            }
            
            color_depth 		= this.color_depth ;
            if(color_depth){
                b= {};
                b[color_depth] = 1;
                color_depth 	= b;
            }else{
                color_depth = null;
            }
            
            language 		= this.fec_lang;
            if(language){
                b= {};
                b[language] = 1;
                language 	= b;
            }else{
                language = null;
            }
            
            login_email_count 		= this.login_email ? 1 : 0;
            register_email_count 	= this.register_email ? 1 : 0;
            
            category_count 		= this.category ? 1 : 0;
            product_count 		= this.sku ? 1 : 0;
            search_count        = 0;
            if (this.search) {
                search_count 		= this.search.text ? 1 : 0;
            }
            
            // first_page
            first_page 		= this.uuid_first_page ? this.uuid_first_page : 0;
            
            // 该处进行了更正，不应该 first_visit_this_url ，而应该使用 uuid_first_page
            is_return = 0;
            
            if(this.ip_first_page == 1){
                ip_count = 1;
            }else{
                ip_count = 0;
            }
            
            uv = 0;
            if(this.uuid_first_page == 1){
                uv = 1;
                is_return = this.is_return ? this.is_return : 0;
            }
            service_date_str = this.service_date_str ? this.service_date_str : null;
            
            cart = this.cart ? this.cart : null;
            order = this.order ? this.order : null;
            
            
            cart_count = 0;
            order_count = 0;
            order_no_count = 0;
            success_order_count = 0;
            success_order_no_count = 0;
            order_amount = 0;
            success_order_amount = 0;
            
            if(cart){
                for(x in cart){
                    one 		= cart[x];
                    if(one && one['qty']){
                        //$sku 		= one['sku'];
                        var skuqty = Number(one['qty'])
                        skuqty = isNaN(skuqty) ? 0 : skuqty
                        cart_count += skuqty;
                    }
                }
            }
            
            if(order && order['invoice'] && order['products']){
                products = order['products'];
                payment_status = order['payment_status'];
                amount = order['amount'];
                if(amount){
                    order_amount = amount;
                }
                order_no_count = 1;
                if(payment_status == 'payment_confirmed'){
                    success_order_no_count = 1;
                    if(amount){
                        success_order_amount = amount;
                    }
                }
                for(x in products){
                    one = products[x];
                    if(one && one['qty']){
                        qty = Number(one['qty']);
                        qty = isNaN(qty) ? 0 : qty
                        order_count += qty;
                        if(payment_status == 'payment_confirmed'){
                            success_order_count += qty;
                            
                        }
                    }
                }
            }
            
            is_return = Number(is_return);
            is_return = isNaN(is_return) ? 0 : is_return
            first_page = Number(first_page);
            first_page = isNaN(first_page) ? 0 : first_page
            emit(service_date_str+"_"+website_id,{
                browser_name:browser_name,
                pv:1,
                uv:uv,
                ip_count:ip_count,
                rate_pv:rate_pv,
                stay_seconds:stay_seconds,
                jump_out_count:jump_out_count,
                drop_out_count:drop_out_count,
                devide:devide,
                country_code:country_code,
                operate:operate,
                fec_app:fec_app,
                is_return: is_return,
                first_page: first_page,
                resolution:resolution,
                color_depth:color_depth,
                language:language,
                login_email_count:login_email_count,
                register_email_count:register_email_count,
                cart_count:cart_count,
                order_count:order_count,
                success_order_count:success_order_count,
                success_order_no_count:success_order_no_count,
                order_no_count:order_no_count,
                order_amount:order_amount,
                success_order_amount:success_order_amount,
                category_count:category_count ,	
                product_count:product_count, 
                search_count:search_count,
                service_date_str:service_date_str
                
            });
            
            
        }
    `
    
    reduceStr := `
        function(key,emits){
            this_pv 				= 0;
            this_rate_pv			= 0;
            this_uv 				= 0;
            this_ip_count 			= 0;
            this_stay_seconds 		= 0;
            this_jump_out_count		= 0;
            this_drop_out_count		= 0 ;
            this_service_date_str 	= null;
            this_devide				= {};
            this_country_code		= {};
            this_browser_name		= {};
            this_operate			= {};
            this_fec_app      		= {};
            this_is_return			= 0;
            this_first_page			= 0;
            this_resolution			= {};
            this_color_depth		= {};
            this_language			= {};
            this_login_email_count	= 0;
            this_register_email_count	= 0;
            this_cart_count				= 0;
            this_order_count			= 0;
            this_success_order_count	= 0;
            this_order_amount			= 0;
            this_success_order_amount	= 0;
            this_success_order_no_count	= 0;
            this_order_no_count	        = 0;
            this_category_count 		= 0;	
            this_product_count			= 0; 
            this_search_count			= 0;
            this_browser_name_2 = [];
            
            for(var i in emits){
                
                if(emits[i].pv){
                    this_pv 			+= emits[i].pv;
                }
                if(emits[i].rate_pv){
                    this_rate_pv 		+= emits[i].rate_pv;
                }
                if(emits[i].uv){
                    this_uv 			+= emits[i].uv;
                }
                if(emits[i].ip_count){
                    this_ip_count 			+= emits[i].ip_count;
                }
                if(emits[i].stay_seconds){
                    this_stay_seconds 	+= emits[i].stay_seconds;
                }
                if(emits[i].service_date_str){
                    this_service_date_str = emits[i].service_date_str;
                }
                //this_customer_id = this_customer_id.concat(emits[i].customer_id);
                if(emits[i].jump_out_count){
                    this_jump_out_count += emits[i].jump_out_count;
                }
                if(emits[i].drop_out_count){
                    this_drop_out_count += emits[i].drop_out_count;
                }
                if(emits[i].devide){
                    devide = emits[i].devide;
                    for(brower_ne in devide){
                        count = devide[brower_ne];
                        if(!this_devide[brower_ne]){
                            this_devide[brower_ne] = count;
                        }else{
                            this_devide[brower_ne] += count;
                        }
                    }
                }
                if(emits[i].country_code){
                    country_code = emits[i].country_code;
                    for(brower_ne in country_code){
                        count = country_code[brower_ne];
                        if(!this_country_code[brower_ne]){
                            this_country_code[brower_ne] = count;
                        }else{
                            this_country_code[brower_ne] += count;
                        }
                    }
                }
                if(emits[i].browser_name){
                    browser_name = emits[i].browser_name;
                    for(brower_ne in browser_name){
                        count = browser_name[brower_ne];
                        if(!this_browser_name[brower_ne]){
                            this_browser_name[brower_ne] = count;
                        }else{
                            this_browser_name[brower_ne] += count;
                        }
                    }
                }
                if(emits[i].operate){
                    operate = emits[i].operate;
                    for(brower_ne in operate){
                        count = operate[brower_ne];
                        if(!this_operate[brower_ne]){
                            this_operate[brower_ne] = count;
                        }else{
                            this_operate[brower_ne] += count;
                        }
                    }
                }
                if(emits[i].fec_app){
                    fec_app = emits[i].fec_app;
                    for(brower_ne in fec_app){
                        
                        count = fec_app[brower_ne];
                        if(!this_fec_app[brower_ne]){
                            this_fec_app[brower_ne] = count;
                        }else{
                            this_fec_app[brower_ne] += count;
                        }
                    }
                }
                if(emits[i].resolution){
                    resolution = emits[i].resolution;
                    for(brower_ne in resolution){
                        count = resolution[brower_ne];
                        if(!this_resolution[brower_ne]){
                            this_resolution[brower_ne] = count;
                        }else{
                            this_resolution[brower_ne] += count;
                        }
                    }
                }
                if(emits[i].color_depth){
                    color_depth = emits[i].color_depth;
                    for(brower_ne in color_depth){
                        
                        count = color_depth[brower_ne];
                        if(!this_color_depth[brower_ne]){
                            this_color_depth[brower_ne] = count;
                        }else{
                            this_color_depth[brower_ne] += count;
                        }
                    }
                }
                if(emits[i].language){
                    language = emits[i].language;
                    for(brower_ne in language){
                        
                        count = language[brower_ne];
                        if(!this_language[brower_ne]){
                            this_language[brower_ne] = count;
                        }else{
                            this_language[brower_ne] += count;
                        }
                    }
                }
                if(emits[i].login_email_count){
                    this_login_email_count 		+= emits[i].login_email_count;
                }
                if(emits[i].register_email_count){
                    this_register_email_count 	+= emits[i].register_email_count;
                }
                if(emits[i].cart_count){
                    this_cart_count 			+= emits[i].cart_count;
                }
                if(emits[i].order_count){
                    this_order_count 			+= emits[i].order_count;
                }
                if(emits[i].success_order_count){
                    this_success_order_count 	+= emits[i].success_order_count;
                }
                if(emits[i].order_amount){
                    this_order_amount 			+= emits[i].order_amount;
                }
                if(emits[i].success_order_amount){
                    this_success_order_amount 	+= emits[i].success_order_amount;
                }
                if(emits[i].success_order_no_count){
                    this_success_order_no_count 	+= emits[i].success_order_no_count;
                }
                if(emits[i].order_no_count){
                    this_order_no_count += emits[i].order_no_count;
                }
                if(emits[i].is_return){
                    this_is_return 			+= emits[i].is_return;
                }
                if(emits[i].first_page){
                    this_first_page 		+= emits[i].first_page;
                }
                if(emits[i].category_count){
                    this_category_count 			+= emits[i].category_count;
                }
                if(emits[i].product_count){
                    this_product_count 			+= emits[i].product_count;
                }
                if(emits[i].search_count){
                    this_search_count 			+= emits[i].search_count;
                }
            }
            
            return {	
                pv:this_pv,
                rate_pv:this_rate_pv,
                uv:this_uv,
                ip_count:this_ip_count,
                jump_out_count:this_jump_out_count,
                drop_out_count:this_drop_out_count,
                stay_seconds:this_stay_seconds, 
                devide:this_devide,
                country_code:this_country_code,
                browser_name:this_browser_name,
                operate:this_operate,
                fec_app:this_fec_app,
                is_return:this_is_return,
                first_page:this_first_page,
                resolution:this_resolution,
                color_depth:this_color_depth,
                language:this_language,
                
                login_email_count:this_login_email_count,
                register_email_count:this_register_email_count,
                cart_count:this_cart_count,
                order_count:this_order_count,
                success_order_count:this_success_order_count,
                order_amount:this_order_amount,
                success_order_amount:this_success_order_amount,
                success_order_no_count:this_success_order_no_count,
                order_no_count:this_order_no_count,
                category_count:this_category_count ,	
                product_count:this_product_count, 
                search_count:this_search_count,
                service_date_str:this_service_date_str
            };
        }
    `
    
    finalizeStr := `
        function (key, reducedVal) {
            uv = reducedVal.uv;
            pv = reducedVal.pv;
            // 平均pv
            if(uv){
                this_pv_rate = pv/uv;
                this_pv_rate = this_pv_rate*10000;
                this_pv_rate = Math.ceil(this_pv_rate);
                this_pv_rate = this_pv_rate/10000;
            }else{
                this_pv_rate = 0;
            }
            
            // 销售转换率
            if(uv){
                this_success_order_no_count = reducedVal.success_order_no_count;
                this_sku_sale_rate = this_success_order_no_count/uv;
                this_sku_sale_rate = this_sku_sale_rate*10000;
                this_sku_sale_rate = Math.ceil(this_sku_sale_rate);
                this_sku_sale_rate = this_sku_sale_rate/10000;
            }else{
                this_sku_sale_rate = 0;
            }
            // 订单支付率
            order_no_count = reducedVal.order_no_count;
            if(order_no_count){
                this_success_order_no_count = reducedVal.success_order_no_count;
                this_order_payment_rate = this_success_order_no_count/order_no_count;
                this_order_payment_rate = this_order_payment_rate*10000;
                this_order_payment_rate = Math.ceil(this_order_payment_rate);
                this_order_payment_rate = this_order_payment_rate/10000;
            }else{
                this_order_payment_rate = 0;
            }
            // 跳出率
            if(uv){
                this_jump_out_count = reducedVal.jump_out_count;
                this_jump_out_rate = this_jump_out_count/uv;
                this_jump_out_rate = this_jump_out_rate*10000;
                this_jump_out_rate = Math.ceil(this_jump_out_rate);
                this_jump_out_rate = this_jump_out_rate/10000;
            }else{
                this_jump_out_rate = 0;
            }
            // 退出率
            if(uv){
                this_drop_out_count = reducedVal.drop_out_count;
                this_drop_out_rate = this_drop_out_count/uv;
                this_drop_out_rate = this_drop_out_rate*10000;
                this_drop_out_rate = Math.ceil(this_drop_out_rate);
                this_drop_out_rate = this_drop_out_rate/10000;
            }else{
                this_drop_out_rate = 0;
            }
            
            if(uv){
                this_is_return = reducedVal.is_return;
                this_is_return_rate = (this_is_return/uv); 
                this_is_return_rate = this_is_return_rate*10000;
                this_is_return_rate = Math.ceil(this_is_return_rate);
                this_is_return_rate = this_is_return_rate/10000;
            }else{
                this_is_return_rate = 0;
            }
            
            reducedVal.is_return_rate= this_is_return_rate;
            this_stay_seconds = reducedVal.stay_seconds;
            //平均停留时间
            rate_pv = reducedVal.rate_pv;
            if(rate_pv){
                stay_seconds_rate = (this_stay_seconds/rate_pv); 
                stay_seconds_rate = stay_seconds_rate*10000;
                stay_seconds_rate = Math.ceil(stay_seconds_rate);
                stay_seconds_rate = stay_seconds_rate/10000;
            }else{
                stay_seconds_rate = 0;
            }
            reducedVal.stay_seconds_rate    = stay_seconds_rate;
            reducedVal.jump_out_rate        = this_jump_out_rate;
            reducedVal.drop_out_rate        = this_drop_out_rate;
            reducedVal.sku_sale_rate        = this_sku_sale_rate;
            reducedVal.order_payment_rate   = this_order_payment_rate;
            reducedVal.pv_rate		        = this_pv_rate;
            reducedVal.website_id           = "` + website_id + `"
            return reducedVal;
        }
    `
    // 结果输出的 mongodb collection
    outDoc := bson.M{"replace": outCollName}
    // 执行mapreduce的job struct
    job := &mgo.MapReduce{
        Map:      mapStr,
        Reduce:   reduceStr,
        Finalize: finalizeStr,
        Out:      outDoc,
    }
    // 开始执行map reduce
    err = mongodb.MDC(dbName, collName, func(coll *mgo.Collection) error {
        _, err := coll.Find(nil).MapReduce(job, nil)
        return err
    })
    if err != nil {
        return err
    }
    
    
    // 上面mongodb maoreduce处理完的数据，需要存储到es中
    // es index 的type name
    esWholeAllTypeName :=  helper.GetEsWholeAllTypeName()
    esIndexName := helper.GetEsIndexNameByType(esWholeAllTypeName)
    // es index 的type mapping
    esWholeAllTypeMapping := helper.GetEsWholeAllTypeMapping()
    // 删除index，如果mapping建立的不正确，可以执行下面的语句删除重建mapping
    //err = esdb.DeleteIndex(esIndexName)
    //if err != nil {
    //    return err
    //}
    // 初始化mapping
    err = esdb.InitMapping(esIndexName, esWholeAllTypeName, esWholeAllTypeMapping)
    if err != nil {
        log.Println("esdb.InitMapping:" + esIndexName + "##" + esWholeAllTypeName + "fail")
        return err
    }
    // 同步mongo数据到ES
    // mongodb中的数据总数
    mCount := 0
    // 得到总数
    err = mongodb.MDC(dbName, outCollName, func(coll *mgo.Collection) error {
        var err error
        mCount, err = coll.Count() 
        return err
    })
    if err != nil {
        return err
    }
    numPerPage := helper.BulkSyncCount
    pageNum := int(math.Ceil(float64(mCount) / float64(numPerPage)))
    for i:=0; i<pageNum; i++ {
        err = mongodb.MDC(dbName, outCollName, func(coll *mgo.Collection) error {
            var err error
            var wholeAlls []model.WholeAll
            coll.Find(nil).Skip(i*pageNum).Limit(numPerPage).All(&wholeAlls)
            log.Println("wholeAlls length:")
            log.Println(len(wholeAlls))
            
            if len(wholeAlls) > 0 {
                // 使用bulk的方式，将数据批量插入到elasticSearch
                bulkRequest, err := esdb.Bulk()
                if err != nil {
                    log.Println("444" + err.Error())
                    return err
                }
                for j:=0; j<len(wholeAlls); j++ {
                    wholeAll := wholeAlls[j]
                    wholeAllValue := wholeAll.Value
                    wholeAllValue.Id = wholeAll.Id_
                    log.Println("888")
                    log.Println(esIndexName)
                    log.Println(esWholeAllTypeName)
                    log.Println(wholeAll.Id_)
                    log.Println(wholeAllValue)
                    req := esdb.BulkUpsertTypeDoc(esIndexName, esWholeAllTypeName, wholeAll.Id_, wholeAllValue)
                    bulkRequest = bulkRequest.Add(req)
                }
                bulkResponse, err := esdb.BulkRequestDo(bulkRequest)
                // bulkResponse, err := bulkRequest.Do()
                if err != nil {
                    log.Println("333" + err.Error())
                    return err
                }
                if bulkResponse != nil {
                    log.Println(bulkResponse)
                }
            }
            return err
        })
        
        if err != nil {
            log.Println("2222" + err.Error())
            return err
        }
    }
    if err != nil {
        log.Println("3333" + err.Error())
    }
    
    
    return err
}

/*
type MapReduce struct {
    Map      string      // Map Javascript function code (required)
    Reduce   string      // Reduce Javascript function code (required)
    Finalize string      // Finalize Javascript function code (optional)
    Out      interface{} // Output collection name or document. If nil, results are inlined into the result parameter.
    Scope    interface{} // Optional global scope for Javascript functions
    Verbose  bool
}
*/

/*
type resultValue struct{
    BrowserNameCount int64 `browser_name`
}

var result []struct { 
    Id string `_id`
    Value resultValue 
} dbName string, collName string, outCollName
*/