 const checkKey = function(rule, value, callback) {
   if (value==='') {
     callback(new Error('不能为空'));
   } else if (value != undefined) {
     var partenHalf = /^\s*$/;
     for (var s in value) {
       if (!value[s].match(/^[A-Za-z0-9_-]+$/)) {
         callback(new Error('请输入字母、数字、‘_-’'));
       }
     }
     callback();
   }
 };
  const checkConfigKey = function(rule, value, callback) {
   if (value==='') {
     callback(new Error('不能为空'));
   } else if (value != undefined) {
     var partenHalf = /^\s*$/;
     for (var s in value) {
       if (!value[s].match(/[\u0000-\u00ff]/g)) {
         callback(new Error('请输入字母、数字或者半角字符'));
       }
     }
     callback();
   }
 };
 const checkWords = function(rule, value, callback) {
   if (!value || value == '') {
     callback(new Error('不能为空'));
   } else {
     callback();
   }
 };
 const checkNumber = function(rule, value, callback) {

   if (value==='') {
     callback(new Error('不能为空'));
   } else if (value != undefined && !Number.isInteger(parseInt(value))) {
     callback(new Error('请输入数字值'));
   } else {
     callback();
   }
 };
 var ruleKeyCustom = {
   validator: checkKey,
   trigger: 'blur',
 };
 var ruleCustom = {
   validator: checkWords,
   trigger: 'blur',
 };
 var ruleNumberCustom = {
   validator: checkNumber,
   trigger: 'blur',
 };
var ruleConfigKeyCustom={
  validator: checkConfigKey,
   trigger: 'blur',
}
 module.exports = {
   ruleKeyCustom: ruleKeyCustom,
   ruleCustom: ruleCustom,
   ruleNumberCustom: ruleNumberCustom,
   ruleConfigKeyCustom: ruleConfigKeyCustom,


 }
