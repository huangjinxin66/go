function isLegalPort(str) {
  var reg = /^\d+$/;
  if (reg.test(str)) {
    if (str > 65535) return false;
    else return true;
  }
  return false;
}
function isDecimal(str) {
  var reg_Decimal = /^\d+\.\d+$/;
  return reg_Decimal.test(str);
}
function isNum(str) {
  var reg_Num = /^\d+$/;
  return reg_Num.test(str);
}
function isIP(str) {
  var reg = /^(\d+)\.(\d+)\.(\d+)\.(\d+)$/;
  return reg.test(str);
}
function checkFileName(str) {
  var reg = /^[^\/]+$/;
  var reg1 = /^(?!.*\\.*$)/;
  return reg.test(str) && reg1.test(str);
}

/**
 * 将秒数换成时分秒格式
 * 作者：龙周峰
 */
function formatSeconds(value) {
  var secondTime = parseInt(value);// 秒
  var minuteTime = 0;// 分
  var hourTime = 0;// 小时
  if (secondTime > 60) {//如果秒数大于60，将秒数转换成整数
    //获取分钟，除以60取整数，得到整数分钟
    minuteTime = parseInt(secondTime / 60);
    //获取秒数，秒数取佘，得到整数秒数
    secondTime = parseInt(secondTime % 60);
    //如果分钟大于60，将分钟转换成小时
    if (minuteTime > 60) {
      //获取小时，获取分钟除以60，得到整数小时
      hourTime = parseInt(minuteTime / 60);
      //获取小时后取佘的分，获取分钟除以60取佘的分
      minuteTime = parseInt(minuteTime % 60);
    }
  }
  var result = "" + parseInt(secondTime) + "s";

  if (minuteTime > 0) {
    result = "" + parseInt(minuteTime) + "min" + result;
  }
  if (hourTime > 0) {
    result = "" + parseInt(hourTime) + "h" + result;
  }
  return result;
}
