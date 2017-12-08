<template>
  <div>
   <H1>Apps应用列表</H1>
    <HandleModal ref="handleModal"></HandleModal><br>
    <AuthModal ref="authModal"></AuthModal>
    
    <Table :columns="tableConf.columns" :data="tableData"></Table>
  </div>
</template>

<script>
  import Api from 'common/api';
  import HandleModal from './handleModal.vue';
  import AuthModal from './userAuth.vue';

  export default {
    name: 'AppList',
    data() {
      return {
        tableConf: {
          columns: [
            {
              title: '业务标记',
              key: 'app',
              align: 'center',
              render: (h, param) => h('div', [
                h('a', {
                  domProps: {
                     href: '#/profiles/'+param.row.app,
                  },
                }, param.row.app),
              ]),
            },
            {
              title: '业务类型',
              key: 'appType',
              align: 'center',
              width: '8%',
              render: (h, param) => h('div', [
                h('span', {
                  domProps: {
                     innerHTML:param.row.appType==0?"服务端":(param.row.appType==9?"客户端":"未知"),
                  },
                }),
              ]),
            },
            {
              title: '业务名称',
              key: 'name',
              align: 'center',
            },
            {
              title: '负责人名称',
              key: 'charger',
              align: 'center',
               width: '8%',
            },
            {
              title: '创建时间',
              key: 'createTime',
              align: 'center',
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
                  attrs:{
                	disabled:  !(param.row.permission == 9||param.row.permission ==99),
                  },
                  on: {
                    click: () => {
                      this.modifyItem(param.row);
                    },
                  },
                }, '修改'),
                 h('Button', {
                  props: {
                    type: 'warning',
                    size: 'small',
                  },
                  style: {
                    marginRight: '5px',
                  },
                  attrs:{
                  	disabled:  !(param.row.permission == 9||param.row.permission ==99),
                   },
                  on: {
                    click: () => {
                      this.authList(param.row.app)
                    },
                  },
                }, '用户列表'),
                h('Button', {
                  props: {
                    type: 'error',
                    size: 'small',
                  },
                     style: {
                    marginRight: '5px',
                  },
                  attrs:{
                    	disabled:  !(param.row.permission == 9||param.row.permission ==99),
                     },
                  on: {
                    click: () => {
                      this.delItem(param.row.app);
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
                  on: {
                    click: () => {
                      window.location.href='#/profiles/'+param.row.app;
                    },
                  },
                }, '进入配置'),
              ]),
            },
          ],
          data: [],
        },
       
      };
    },
    components: {
      HandleModal,
      AuthModal
    },
    computed: {
      tableData() {
        return this.$store.state.MOD.tableList;
      },
    },
    mounted() {
      const self = this;    
       self.$store.dispatch('getTableList');
         
      
    },
    methods: {
      authList(app){
         this.$store.state.MOD.currentApp=app;
         this.$refs.authModal.getUserList();
      },
      modifyItem(data) {
        this.$refs.handleModal.setModifyData(data);
      },
      delItem(id) {

        this.$Modal.confirm({
                      title: '提示',
                      content: '<p>确认删除？</p>',
                      onOk: () => {
                          Api.delApp(id).done((res) =>{
                            if(res.code==0){
                              // alert("删除成功");
                              this.$Message.success('删除成功');
                              this.$store.dispatch('getTableList');  
                            }else{
                               this.$Message.error("删除失败："+res.msg);
                            }
                          });
                      },
                      onCancel: () => {
                         
                      }
                  });

       
      }
      
    },
  };
</script>

<style lang="scss">
</style>
