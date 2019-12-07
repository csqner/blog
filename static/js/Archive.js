// JavaScript Document

function iniParam() {
    //初始化WOW.js
    new WOW().init();

    //页面效果
    setTimeout(function () {
        $('#menu-cb').click();
        setTimeout(function () {  $('#menu-cb').click();}, 1500);
    }, 1000);
}