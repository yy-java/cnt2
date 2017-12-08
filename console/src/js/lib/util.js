var getUrlVar = function(key) {
  return decodeURIComponent((new RegExp('[?|&]' + key + '=' + '([^&;]+?)(&|#|;|$)').exec(location.href) || [, ""])[1].replace(/\+/g, '%20')) || null;
};

var getCookie = function(name) {
  var arr,
    RE = new RegExp("(^| )" + name + "=([^;]*)(;|$)");
  if (arr = document.cookie.match(RE)) {
    return unescape(arr[2]);
  } else {
    return ''
  }
};
var setCookie = function (key, value, options) {
  if (arguments.length > 1 && !$.isFunction(value)) {
    options = $.extend({}, options);

    if (typeof options.expires === 'number') {
      var days = options.expires, t = options.expires = new Date();
      t.setMilliseconds(t.getMilliseconds() + days * 864e+5);
    }
    return (document.cookie = [
      encodeURIComponent(key), '=', encodeURIComponent(value),
      options.expires ? '; expires=' + options.expires.toUTCString() : '', // use expires attribute, max-age is not supported by IE
      options.path    ? '; path=' + options.path : '',
      options.domain  ? '; domain=' + options.domain : '',
      options.secure  ? '; secure' : ''
    ].join(''));
  }
};
var clearCookie = function () {
      document.cookie="username=;expires=0;domain=xx.com;";
};
var isEmbed = (function() {

  var ret = false;
  try {
    window.external.sendCommand("getMyUid");
    ret = true;
  } catch (e) {
    ret = false;
  }
  return ret;

})();

var formatNum = function(num) {
  return num.toString().replace(/\d+?(?=(?:\d{3})+$)/img, "$&,");
};

var format = function(text, obj) {
    var c = obj || {};
    var b = text;
    for (var a in c) {
        if (c.hasOwnProperty(a)) {
            b = b.replace(new RegExp("\\$\\{" + a + "\\}", "gm"), c[a]);
        }
    }
    return b;
}

var attachCss3 = function(boxSelector) {
  if (window.PIE) {
    $(boxSelector).find(".rounded").each(function() {
      try {
        PIE.attach(this);
      } catch (e) {}
    });
  }
};

var detachCss3 = function(boxSelector) {
  if (window.PIE) {
    $(boxSelector).find(".rounded").each(function() {
      try {
        PIE.detach(this);
      } catch (e) {}
    });
  }
};

var isMob = (function(W) {
  var ua = W.navigator.userAgent.toLowerCase();
  if (/iphone|ios|android|mobile/i.test(ua)) {
    return true;
  }
})(window);

module.exports = {
  getUrlVar: getUrlVar,
  getCookie: getCookie,
  setCookie: setCookie,
  clearCookie: clearCookie,
  isEmbed: isEmbed,
  isMob: isMob,
  formatNum: formatNum,
  attachCss3: attachCss3,
  detachCss3: detachCss3,
  format: format,
};
