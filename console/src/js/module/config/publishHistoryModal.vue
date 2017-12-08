<template>
  <div>
    <Modal
      :width="783"
      :height="400"
      v-model= "modal"
      :title= "title"
      :mask-closable="true"
      :footer-hide="true"
      class-name="vertical-center-modal"
    >
     <p slot="header" style="color:#f60;text-align:left">
          
            <span>{{title}} </span>
    </p>
     <div style="height:500px;overflow:auto">
     	<Table :columns="tableConf.columns" :data="tableData"></Table>
	 </div>
    </Modal>
  </div>
</template>

<script>
//config.app + "/" + config.profile+"/" + config.key
import Api from 'common/api';
var publishStatus=['未发布','成功','失败'];

export default {
    name: 'publishHistoryList',
    data() {
      return {
      	modal:false,
      	loading:true,
      	config:[],
      	title:'发布历史',
        tableConf: {
          data:[],
          columns: [
           {
              title: 'IP',
              key: 'sip',
              align: 'left',
              width: '20%',
           },
           {
              title: '机房',
              align: 'left',
              width: '25%',
              render: (h, param) => h('div', [
                h('div', {
                  domProps: {
                     innerHTML: "<span class='lookup_ip' data='"+param.row.sip+"'></span>",
                  },
                }),
              ]),
           },
           {
              title: '发布类型',
              align: 'left',
              width: '12%',
              render: (h, param) => h('div', [
                h('span', {
                  domProps: {
                     innerHTML:  param.row.publishType == 1?"全量":"灰度",
                  },
                }),
              ]),
            },
      	    {
              title: '发布状态',
              align: 'left',
              width: '12%',
              render: (h, param) => h('div', [
                h('div', {
                  domProps: {
                     innerHTML: publishStatus[param.row.publishResult],
                  },
                }),
              ]),
            },
            {
              title: '发布时间',
              key: 'publishTime',
              align: 'left',
            },
          ]
        },
        location:{},
        curVersion:''
      };
    },
    computed: {
      tableData() {
        return this.tableConf.data;
      },
    },
    mounted() {
    
    },
    methods: {
      setPublishHistoryData(data) {
         const jsonData = JSON.parse(JSON.stringify(data));
         this.config = jsonData;
         this.modal = true;
         this.curVersion=jsonData.version;
         this.title='发布记录：'+ jsonData.app+" > "+jsonData.profile+" > "+jsonData.key +" > V" + jsonData.curVersion;
         //this.$store.dispatch('getConfigHistoryList',this.config);
          const self = this;
         Api.queryPublishHistory(this.config).done((res) => {
         	if(res && res.data){
	        	self.tableConf.data = res.data;
	        }else{
	        	self.tableConf.data = [];
	        }
	        self.$nextTick(function(){
		       $('.lookup_ip').each(function(){
					var t =$(this);
					Api.queryIp(t.attr("data")).done((resp) => {
						if(resp && resp.data){
							t.html(resp.data);
						}else{
							t.html("--")
						}
					});
		       }); 
	     	});
	     });
	     
      },
    },
  };
</script>

<style lang="scss">
	.vertical-center-modal{
        display: flex;
        align-items: center;
        justify-content: center;
        .ivu-modal{
            top: 0;
        }
    }
</style>
