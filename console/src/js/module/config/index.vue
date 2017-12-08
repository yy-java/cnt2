<template>
  <div>
     <H1>Configs配置</H1>
     <div style="background:#eee;padding: 10px">
       <h3>
        <Tooltip content="应用列表" placement="top-start">
            <a href="#/applist">首页</a>   >
        </Tooltip>
        <Tooltip content="profiles" placement="top-start">
             <a :href="'#/profiles/'+location.app">{{location.app}}</a>
        </Tooltip>
           > {{location.profile}}     
       </h3>
    </div>
        <br>
        <HandleModal ref="handleModal"></HandleModal>
        <ClientButton ref="client"></ClientButton>
        <ConfigHistoryModal ref="configHistoryModal"></ConfigHistoryModal>
        <Publish ref="publish"></Publish>
       <br>
    <ButtonGroup size='large'>
		<Button :type="btnTypeAll" @click="renderAll">全部</Button>
	    <Button :type="btnTypeApprove" @click="waitApprove">待审核</Button>
	    <Button :type="btnTypePublish" @click="waitPublish">待发布</Button>
    </ButtonGroup>
    
    <Table :columns="tableConf.columns" :data="tableData"></Table>
  </div>
</template>

<script>
  import Api from 'common/api';
  import HandleModal from './handleModal.vue';
  import ConfigHistoryModal from './configHistoryModal.vue';
  import ClientButton from './client.vue';
  import Publish from './publish.vue';
  let allDatas = [];
  export default {
    name: 'configList',
    data() {
      return {
      	btnTypeAll:'primary',
      	btnTypeApprove:'ghost',
      	btnTypePublish:'ghost',
      	permission:-1,
        tableConf: {
          columns: [
            {
              title: 'Key',
              key: 'key',
              align: 'left',
              width: '12%',
            },
             {
              title: '描述',
              key: 'description',
              align: 'left',
              width: '12%',
            },
            {
              title: 'Value',
              key: 'value',
              align: 'left',
              width: '55%',
              render: (h, param) => h('div', [
                h('div', {
                  domProps: {
                     innerHTML: '<textarea class="node-value" disabled>'+param.row.value+'</textarea>',
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
                  },
                  on: {
                    click: () => {
                      this.modifyItem(param.row);
                    },
                  },
                }, '修改'),
                h('Button', {
                  props: {
                    type: 'error',
                    size: 'small',
                  },
                   style: {
                    marginRight: '5px',
                  },
                  attrs:{
                  	disabled: !(param.row.publishedVersion == 0 || this.permission == 9),//刚创建的还未发布的，是可以删的。管理员也可以删
                  },
                  on: {
                    click: () => {
                      this.delItem(param.row);
                    },
                  },
                }, '删除'),
                h('Button', {
                  props: {
                    type: 'success',
                    size: 'small',
                  },
                  style: {
                    marginRight: '5px',
                  },
                  attrs:{
                  	disabled: !window.localStorage, 
                  	title:!window.localStorage?"当前浏览器不支持":"",
                  },
                  on: {
                    click: () => {
                      this.copyData(param.row);
                    },
                  },
                }, '复制'),
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
						this.listConfigHistory(param.row);
                    },
                  },
                }, '历史'),
                
                h('Button', {
                  props: {
                    type: 'info',
                    size: 'small',
                  },
                  attrs:{
                  	disabled: !(param.row.approveType ==1 && this.permission == 9),
                  },
                  style:{
                    marginRight: '5px',
                  },
                  
                  on: {
                    click: () => {
						this.approve(param.row);
                    },
                  },
                }, '审核'),
                h('Button', {
                  props: {
                    type: 'info',
                    size: 'small',
                  },
                  attrs:{
                  	disabled: !(param.row.approveType==2 && param.row.version > param.row.publishedVersion),
                  },
                  style:{
                    marginRight: '5px',
                  },
                  on: {
                    click: () => {
						this.publish(param.row);
                    },
                  },
                }, '发布'),
              ]),
            },
          ],
          data: [],
        },
        location:{},
        status:0//默认显示全部，1是待审核，2是待发布
      };
    },
    computed: {
      tableData(){
        if(this.status==0&&this.$store.state.MOD.configKVList){
          return this.$store.state.MOD.configKVList;
        }else if(this.status==1&&this.$store.state.MOD.waitApproveConfig){
          return this.$store.state.MOD.waitApproveConfig;
        }else if(this.status==2&&this.$store.state.MOD.waitPublishConfig){
            return this.$store.state.MOD.waitPublishConfig;
        }
      }
    },
    components: {
      HandleModal,
      ConfigHistoryModal,
       ClientButton,
       Publish
    },
   
    mounted() {
      const self = this;
      self.location.app=this.$route.params.app;//重复设置
      self.location.profile=this.$route.params.profile;
      self.$store.state.MOD.currentApp=self.location.app;
      self.$store.state.MOD.currentProfile=self.location.profile;
      //每次进入页面先清除旧数据
      self.$store.state.MOD.configKVList=[];
      self.$store.state.MOD.waitApproveConfig=[];
      self.$store.state.MOD.waitPublishConfig=[];
     
       //查询当前用户的权限
      Api.queryCurUserAuth(this.$route.params.app).done((res) => {
          if(res && res.code == 0){
          	self.permission = res.data == 99?9 : res.data;
          }
          self.$store.dispatch('getConfigList');
      });
     
    },
    methods: {
      modifyItem(data) {
         this.$refs.handleModal.setModifyData(data);
      },
      copyData(data) {
        const json =  JSON.stringify(data) ;
        window.localStorage.setItem(data.key,json);
        window.localStorage.setItem("cnt2-copy-key",data.key);
        this.$Message.success('已复制到粘贴板');
      },
      delItem(data) {
        const self = this;
         this.$Modal.confirm({
                    title: '提示',
                    content: '<p>确认删除？</p>',
                    onOk: () => {
                        Api.delConfig(data).done((res) => {
                          if(res.code==0){
                            this.$Message.success('删除成功');
                            self.$store.dispatch('getConfigList');
                          }else{
                             this.$Message.error("删除失败："+res.msg);
                          }
                        });
                    },
                    onCancel: () => {
                       
                    }
                });

       
      },
      listConfigHistory(data) {
		 this.$refs.configHistoryModal.setConfigHistoryData(data);
      },
      waitApprove(){
      	// this.tableConf.data=this.$store.state.MOD.waitApproveConfig;
         this.status=1;
      	 this.btnTypeAll  =  'ghost';
      	 this.btnTypePublish = 'ghost';
      	 this.btnTypeApprove  =  'primary';
      },
      waitPublish( ){
      	 // this.tableConf.data=this.$store.state.MOD.waitPublishConfig;
          this.status=2;
      	 this.btnTypeAll  =  'ghost';
      	 this.btnTypeApprove= 'ghost';
      	 this.btnTypePublish  =  'primary';
      },
     renderAll(){
         this.status=0;
      	 // this.tableConf.data = this.$store.state.MOD.configKVList;
         this.btnTypeAll  =  'primary';
         this.btnTypeApprove= 'ghost';
      	 this.btnTypePublish  =  'ghost';
      },
      publish(data){
         this.$refs.publish.queryNodeList(data);
     //  	 if(!data || !data.id){
     //  	 	return;
     //  	 }
     //  	 Api.publish(data.id,nodes).done((res) => {
		  	// if(res && res.code == 0){
		  	// 	this.$Message.success('发布成功');
		  	// 	location.reload();
		  	// }else{
		  	// 	this.$Message.error(res?res.msg:'发布错误');
		  	// }
     //  	});
      },
      approve(data){
      	 if(!data || !data.app){
      	 	return;
      	 }
      	 Api.approve({'app':data.app,'profile':data.profile,'key':data.key,'version':data.version}).done((res) => {
		  	if(res && res.code == 0){
		  		this.$Message.success('审核成功');
		  		this.$store.dispatch('getConfigList');
		  	}else{
		  		this.$Message.error(res?res.msg:'审核错误');
		  	}
      	});
      },
    },
  };
</script>

<style lang="scss">
.node-value{width:100%; height:100%; min-height:80px; min-width: 500px;}
textarea{background-color: #f8f8f9;}
</style>
