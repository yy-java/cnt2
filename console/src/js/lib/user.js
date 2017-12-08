var $ = require('jquery');
var Util = require('./util');


var isLogin = function() {
  return !!Util.getCookie('username');
};

var showLoginBox = function(url) {
  var _url = url || location.href;
  UDBSDKProxy.openByFixProxy(_url);
};

var getUid = function() {
  return Util.getCookie('uid') || 0;
};
var getUsername = function() {
  return Util.getCookie("username");
};

module.exports = {
  // createLoginIframe: createLoginIframe,
  isLogin: isLogin,
  showLoginBox: showLoginBox,
  getUid: getUid,
  getUsername: getUsername,
  isYyStaff: function() {
    return true;
    // return Util.getCookie('username').indexOf('dw_') === 0;
  }
};
