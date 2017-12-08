<template>
  <div>
    <Button type="success" @click="type=1;modal = true;formData={
          app: '',
          name: '',
          appType: 0,
          charger: getName(),
          chargerUid:getUid()
      }">新增APP
    </Button>
    <Modal
      :width="800"
      v-model="modal"
      :title="type==1?'新增APP':'修改APP'"
      ok-text="提交"
      :mask-closable="false"
      :footer-hide="true"
    >
      <Form :model="formData" ref="formData">
        <Form-item :rules="ruleKeyCustom" prop="app" label="业务标记：" required>
          <Input v-model="formData.app"  :readonly="type==1?false:true" :disabled="type==1?false:true">
          </Input>
        </Form-item>
        <Form-item :rules="ruleCustom" prop="name" label="业务名称：" required>
          <Input v-model="formData.name" >
          </Input>
        </Form-item>
        <Form-item  prop="appType" label="业务类型：" required>
          <RadioGroup v-model="formData.appType" >
              <Radio label="0" :disabled="type==1?false:true">
                  <span>服务端</span>
              </Radio>
              <Radio label="9" :disabled="type==1?false:true">
                  <span>客户端</span>
              </Radio>
          </RadioGroup>
        </Form-item>
        <Form-item :rules="ruleCustom" prop="charger" label="负责人：" required>
          <Input v-model="formData.charger" :disabled="type==1?false:true">
          </Input>
        </Form-item>

        <Form-item :rules="ruleNumberCustom" prop="chargerUid" label="负责人UID：" required>
          <Input v-model="formData.chargerUid" :disabled="type==1?false:true">
          </Input>
        </Form-item>
        <Form-item>
          <Button type="primary" @click="submitData('formData')">
            {{type === 1 ? '提交' : '修改'}}
          </Button>
        </Form-item>
      </Form>

    </Modal>
  </div>
</template>

<script>
  import Api from 'common/api';
  import User from 'lib/user';
  import Check from 'common/check';
  export default {
    name: 'handleModal',
    data() {
      return {
        modal: false,
        loading: true,
        type: 1, // 1:新增 2：修改
        formData: {
          app: '',
          name: '',
          appType: 0,
          chargerUid: '',
          charger: '',
        },
        ruleCustom:  Check.ruleCustom,
        ruleNumberCustom: Check.ruleNumberCustom, 
        ruleKeyCustom: Check.ruleKeyCustom,  
      };
    },
    updated:function(){
      if(this.type==2){
         this.$refs['formData'].validate();
      }
    },
    methods: {
      getUid(){
        return User.getUid();
      },
     getName(){
      return User.getUsername();
     },
      setModifyData(data) {
       
        const formData = JSON.parse(JSON.stringify(data));
        this.formData = formData;
        this.modal = true;
        this.type = 2;
      },
      submitData(name) {
        // var keys=['app','name','charger','chargerUid'];
        // if(!Check.validateFormData(this.formData,keys)){
        //     this.$Message.error('选项不能为空');
        //     return;
        // }
        this.$refs[name].validate((valid) => {
            if (valid) {
                const self = this;
                if (self.type === 1) {
                  Api.createApp(self.formData).done((res) => {
                    if(res.code === 0){
                      this.$Message.success('添加成功');
                      this.$store.dispatch('getTableList');  
                      this.modal = false;
                    }else{
                        this.$Message.error('添加失败，' + res.msg);
                    }
                  });
                } else {
                  Api.updateApp(self.formData).done((res) => {
                    if(res.code === 0){
                      this.$Message.success('修改成功');
                      this.$store.dispatch('getTableList');  
                       this.modal = false;
                    }else{
                      this.$Message.error('修改失败，' + res.msg);
                    }
                  });
                }
            } else {
                this.$Message.error('表单验证失败!');
            }
        });

      },
    },
    
  };
</script>

<style lang="scss">
</style>
