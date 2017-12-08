<template>
 	<span>
  	<Button type="info" @click="queryNodeList">在线客户端
  	</Button>
    <Modal
      :width="1000"
      v-model="modal" 
      :mask-closable="false"
      :footer-hide="true"
    >  
    <p slot="header" style="color:#f60;text-align:center">
          
            <span>在线客户端列表</span>
    </p>
    <Table :columns="tableUser.columns" :data="tableUser.data"></Table>
    </Modal>
  </span>
</template>
<script>
  import Api from 'common/api';
  
  export default {
    name: 'client',
    data() {
      return {
        tableUser:{
          columns:[ 
           {
              title: 'nodeId',
              key: 'nodeId',
              align: 'center',
            },
            {
              title: 'pid',
              key: 'pid',
              align: 'center',
            },
            {
              title: 'ip',
              key: 'sip',
              align: 'center',
            },
            {
              title: '注册时间',
              key: 'registerTime',
              align: 'center',
              render: (h, param) => h('div', [
                h('span', {
                  domProps: {
                     innerHTML: new Date(param.row.registerTime).Format("yyyy-MM-dd hh:mm:ss"),
                  },
                }),
              ]),
            },
           
            ],
            data:[]
        },
        modal: false,
        loading: true,
      };
    },    
  	methods:{
  		queryNodeList(){
  			this.modal=true;			
  			let formData={         
		         app:this.$store.state.MOD.currentApp, 
		      	 profile:this.$store.state.MOD.currentProfile,
        	};
  			Api.getNodeList(formData).done((res)=>{ 
			 	if(res && res.data && res.data.nodes && res.data.nodes.length){
    	 			this.tableUser.data = res.data.nodes;    
			 	}else{
			 		this.tableUser.data = [];
			 	}
      		});
  		} 	
   }
  };
</script>
<style lang="scss">
</style>
