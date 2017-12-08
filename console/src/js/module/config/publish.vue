<template>
 	<span>
    <Modal
      :width="1000"
      v-model="modal" 
      :mask-closable="false"
      :footer-hide="true"
    >  
    <p slot="header" style="color:#f60;text-align:left">
          
            <span>{{title}}</span>
    </p>
    <Button type="primary" @click='publishNode'> &nbsp;&nbsp;发&nbsp;&nbsp;&nbsp;&nbsp;布&nbsp;&nbsp;</Button>     &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
    <Button type="error" @click='publishAll'>全量发布</Button><br><br> 
    <Table  border :columns="tableNode.columns" :data="tableNode.data" @on-select="onSelect" @on-select-all="onSelectAll" @on-selection-change="onSelectChange"></Table>
    </Modal>
  </span>
</template>
<script>
  import Api from 'common/api';
  
  export default {
    name: 'publish',
    data() {
      return {
        tableNode:{
          columns:[ 
           {
              type: 'selection',
              width: 60,
              align: 'center',
           },
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
            {
              title: '状态',
              align: 'center',
              render: (h, param) => h('div', [
                h('span', {
                  domProps: {
                     innerHTML:param.row._disabled?"已发布":"未发布",
                  },
                }),
              ]),
            },
           
            ],
            data:[]
        },
        title:'',
        modal: false,
        loading: true,
        selectedRow:[],
        allUnPublishedNodes:[],
        configId:''
      };
    }, 
  	methods:{
      onSelect(selection,row){
        this.selectedRow=selection;
      },
      onSelectAll(selection){
      },
      onSelectChange(selection){
        this.selectedRow=selection;
      },
      publishNode(event,isAll){
      	let data={
	          id:this.configId,
	          nodes:[],
	          pubType:1,
	        };
	    if(isAll){
	    	data.nodes=this.allUnPublishedNodes;
	    }else{
	        //全选不需要传nodes参数
	        if(!this.selectedRow||this.selectedRow.length==0){
	          this.$Message.info("请选择节点");
	          return;
	        }
	        for(var index in this.selectedRow){
	           data.nodes.push(this.selectedRow[index].nodeId);
	        }
	    }
	    if(data.nodes.length == this.allUnPublishedNodes.length){
	    	data.pubType=2;
	    }  
        let self=this;
        data.nodes=data.nodes.join();
        Api.publish(data).done((res)=>{
           if(res.code==0){
            self.$Message.success("发布成功");
            self.$store.dispatch('getConfigList');
            this.modal=false;
           }else{
            self.$Message.error("发布失败："+res.msg);

           }
        });
      },
      publishAll(){
      	 this.publishNode(null,true);
      },
  	 queryNodeList(data){
		this.configId=data.id;
		this.modal=true;
		this.selectedRow = [];
		this.allUnPublishedNodes = [];	
		let formData={         
	         app:data.app, 
	      	 profile:data.profile,
	      	 key:data.key,
	      	 version:data.version,
    	};
        let self=this;
        this.title='节点发布：'+ data.app+" > "+data.profile+" > "+data.key;
  		Api.getNodeList(formData).done((res)=>{
    			 	if(res && res.data && res.data.nodes && res.data.nodes.length){
    			 		let nodes = res.data.nodes;
    			 		let publishedNodes = res.data.publishedNodes;
    			 		
			 			for(var i in nodes){
			 				nodes[i]["_disabled"] = false
			 				for(var j in publishedNodes){
			 					if(nodes[i].nodeId == publishedNodes[j]){
				 					nodes[i]["_disabled"] = true
				 					break;
			 					}
			 				}
			 				if(nodes[i]["_disabled"] == false){
			 					self.allUnPublishedNodes.push(nodes[i].nodeId)
			 				}
			 			}
        	 			self.tableNode.data = nodes;    
    			 	}else{
    			 		self.tableNode.data= [];
    			 	}
      		});
  		} 	
   }
  };
</script>

<style lang="scss">
</style>
