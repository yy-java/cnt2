<template>
  <span>
   
    <Button type="success" @click="type=1;modal = true;formData={        
        key: '',
        value: '',
        validator: '',
        description: '',
      }">添加配置
    </Button>
   
    <Modal
      :width="800"
      v-model="modal"
      :title="type==1?'新增配置':'修改配置'"
      ok-text="提交"
      :mask-closable="false"
      :footer-hide="true"
    >
      <Form :model="formData" ref="formData">
        <Form-item :rules="ruleKeyCustom" prop="key" label="key：" required>
          <Input v-model="formData.key" :readonly="type==1?false:true" :disabled="type==1?false:true">
          </Input>
        </Form-item>
        <Form-item :rules="ruleCustom" prop="value" label="value：" required>
          <Input v-model="formData.value" type="textarea" :rows="4">
          </Input>
        </Form-item>
       <Form-item  prop="validator" label="validator：">
          <Input v-model="formData.validator"  type="textarea">
          </Input>
        </Form-item>
         <Form-item :rules="ruleCustom" prop="description" label="描述：" required>
          <Input v-model="formData.description" type="textarea"  >
          </Input>
        </Form-item>

        <Form-item>
          <Button type="primary"  @click="submitData('formData')">
            {{type === 1 ? '保 存' : '修 改'}}
          </Button>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
          <Button type="success" :disabled="canCopy?false:true"  :title="canCopy?'':'当前浏览器不支持'" @click="copyData()">
           复制 
          </Button>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
          <Button type="warning" :disabled="canCopy?false:true"  :title="canCopy?'':'当前浏览器不支持'" @click="pasteData()">
           粘贴
          </Button>
         
        </Form-item>
      </Form>

    </Modal>
  </span>
</template>

<script>
  import Api from 'common/api';
  import Check from 'common/check';
  export default {
    name: 'handleModal',
    data() {
      return {
        modal: false,
        loading: true,
        type: 1, // 1:新增 2：修改
        canCopy: !!window.localStorage ,
        formData: {   
          // app:this.$route.params.app ,
          // profile:this.$route.params.profile,
          key: '',
          value: '',
          validator: '',
          description: ''
        },
        ruleCustom:  Check.ruleCustom,
        ruleKeyCustom: Check.ruleConfigKeyCustom,  
      };
    },
     updated:function(){
      if(this.type==2){
         this.$refs['formData'].validate();
      }
    },
    methods: {
      setModifyData(data) {
        const formData = JSON.parse(JSON.stringify(data));
        this.formData = formData;
        this.modal = true;
        this.type = 2;
      },
      copyData( ) {
        const json =  JSON.stringify(this.formData) ;
        window.localStorage.setItem(this.formData.key,json);
        window.localStorage.setItem("cnt2-copy-key",this.formData.key);
        this.$Message.success('已复制到粘贴板');
      },
      pasteData( ) {
      	const copyedKey =  window.localStorage.getItem("cnt2-copy-key");
      	const jsonObj = JSON.parse(window.localStorage.getItem(copyedKey));
      	if(this.type == 1 || this.formData.key == jsonObj.key){
      		 this.formData = jsonObj;
      	} 
      },
      submitData(name) {
        const self = this;
         this.$refs[name].validate((valid) => {
            if (valid) {
                self.formData.app=this.$route.params.app ;
                self.formData.profile=this.$route.params.profile;

                if (self.type === 1) {     
                    Api.createConfig(self.formData).done((res) => {
                    if(res.code === 0){
                      	this.$Message.success('添加成功');
                        self.$store.dispatch('getConfigList');
                        this.modal=false;
                    }else{
                     	this.$Message.error('添加失败，' + res.msg);
                    }
                  });
                } else {
                  Api.updateConfig(self.formData).done((res) => {
                    if(res.code === 0){
                      this.$Message.success('修改成功');
                      self.$store.dispatch('getConfigList');
                       this.modal=false;
                    }else{
                       this.$Message.error('修改失败，' + res.msg);
                    }
                  });
                }
              }else {
                this.$Message.error('表单验证失败!');
            }
        });
      },
    },
  };
</script>

<style lang="scss">
</style>
