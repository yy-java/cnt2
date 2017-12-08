<template>
  <div>
    <Modal
      :width="1024"
      :height="600"
      v-model= "modal"
      :title= "title"
      :mask-closable="true"
      :footer-hide="true"
      :scrollable=true
      class-name="vertical-center-modal"
    >
     <p slot="header" style="color:#f60;text-align:left">
          
            <span>{{title}} </span>
    </p>
    <div style="height:600px;overflow:auto">
	   <Table :columns="tableConf.columns" :data="tableData"></Table>
    </div>
    
    </Modal>
    <br>
     <PublishHistoryModal ref="publishHistoryModal"></PublishHistoryModal>
  </div>
  
</template>

<script>
//config.app + "/" + config.profile+"/" + config.key
import Api from 'common/api';
import PublishHistoryModal from './publishHistoryModal.vue';
var operateTypes=['新增','修改','修改','回滚'];
export default {
    name: 'configHistoryList',
    data() {
      return {
      	modal:false,
      	loading:true,
      	config:[],
      	title:'配置历史',
        tableConf: {
          data:[],
          columns: [
      	    {
              title: '版本',
              key: 'curVersion',
              align: 'left',
              width: '6%',
              render: (h, param) => h('div', [
                h('div', {
                  domProps: {
                     innerHTML: 'V' + param.row.curVersion,
                  },
                }),
              ]),
            },
            {
              title: '修改人',
              key: 'modifier',
              align: 'left',
              width: '15%',
            },
            {
              title: '修改时间',
              key: 'modifyTime',
              align: 'left',
              width: '15%',
            },
            {
              title: '操作类型',
              key: 'operateType',
              align: 'left',
              width: '10%',
              render: (h, param) => h('div', [
                h('span', {
                  domProps: {
                     innerHTML: operateTypes[param.row.operateType-1],
                  },
                }),
              ]),
            },
            {
              title: 'Value',
              key: 'value',
              align: 'left',
              width: '35%',
              render: (h, param) => h('div', [
                h('div', {
                  domProps: {
                     innerHTML: '<textarea class="node-value" disabled>'+param.row.curValue+'</textarea>',
                  },
                }),
              ]),
            },
            {
              title: '操作',
              key: 'action',
              align: 'center',
              render: (h, param) => h('div', [
                h('Button', {
                  props: {
                    type: 'primary',
                    size: 'small',
                  },
                  style: {
                    marginRight: '5px',
                     //display: (this.curVersion == param.row.curVersion)?"none":"all",
                  },
                  attrs:{
                  	disabled:  (this.curVersion == param.row.curVersion),
                  },
                  on: {
                    click: () => {
                      this.rollback(param.row);
                    },
                  },
                }, '回滚'),
                h('Button', {
                  props: {
                    type: 'warning',
                    size: 'small',
                  },
                  style: {
                    marginRight: '5px',
                  },
                  on: {
                    click: () => {
                      this.publishHistory(param.row);
                    },
                  },
                }, '发布记录'),
              ]),
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
   components: {
      PublishHistoryModal,
    },
    methods: {
      publishHistory(data) {
		 this.$refs.publishHistoryModal.setPublishHistoryData(data);
      },
      setConfigHistoryData(data) {
         const jsonData = JSON.parse(JSON.stringify(data));
         this.config = jsonData;
         this.modal = true;
         this.curVersion=jsonData.version;
         this.title='配置历史：'+ jsonData.app+" > "+jsonData.profile+" > "+jsonData.key;
         //this.$store.dispatch('getConfigHistoryList',this.config);
          const self = this;
         Api.queryConfigHistory(this.config).done((res) => {
	        self.tableConf.data = res.data;
	     });
      },
      rollback(data) {
        const self = this;
		 Api.rollbackConfig(data).done((res) => {
		 	if(!res){
		 		self.$Message.error('请求异常')
		 		return;
		 	}
		    if(res.code == 0){
		    	self.$Message.success("操作成功！")
		 		   self.$store.dispatch('getConfigList');
           self.modal=false;
		 	}else{
		 		alert(res.msg);
		 	}
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