"use strict";var _typeof="function"==typeof Symbol&&"symbol"==typeof Symbol.iterator?function(e){return typeof e}:function(e){return e&&"function"==typeof Symbol&&e.constructor===Symbol&&e!==Symbol.prototype?"symbol":typeof e};!function(e){"function"==typeof define&&define.amd?define(["jquery"],e):e("object"===("undefined"==typeof exports?"undefined":_typeof(exports))?require("jquery"):jQuery)}(function(e){function t(e){var t=/^.*((youtu.be\/)|(v\/)|(\/u\/\w\/)|(embed\/)|(watch\?))\??v?=?([^#\&\?]*).*/,i=e.match(t);return i&&11==i[7].length?i[7]:!1}function i(e){var t=/(\/id_)(\w+)/,i=/http?:\/\/(www\.)|(v\.)youku.com/,o=e.match(t);return i.test(e)&&o&&o[2]?o[2]:!1}function o(t,i){this.$element=e(t),this.options=e.extend({},o.DEFAULTS,e.isPlainObject(i)&&i),this.init()}var r="qor.medialibrary.action",a="enable."+r,n="disable."+r,d="keyup."+r,l="switched.qor.tabbar.radio",s='[data-toggle="qor.tab.radio"]',f="reload.qor.bottomsheets",c='[name="QorResource.SelectedType"]',u='[data-tab-source="video_link"]',m=".qor-video__link",h=".qor-medialibrary__video-link",b=".qor-medialibrary__video",p=".qor-medialibrary__desc",v=".qor-file__options",y=".qor-bottomsheets__mediabox",g='input[name="QorResource.MediaOption"]';return o.prototype={constructor:o,init:function(){this.bind(),this.initMedia()},bind:function(){e(document).on(l,s,this.resetMediaData.bind(this)).on(d,m,this.setVideo.bind(this)).on(d,p,this.setImageDesc.bind(this)).on(f,y,this.initMedia.bind(this,"bottomsheet"))},unbind:function(){e(document).off(l,s,this.resetMediaData.bind(this)).off(d,m,this.setVideo.bind(this)).off(d,p,this.setImageDesc.bind(this))},setMediaData:function(e,t){var i=e.find(v),o=e.find(g);i.val(JSON.stringify(t)),o.val(JSON.stringify(t))},setImageDesc:function(t){var i,o,r=e(t.target),a=r.closest("form");i=a.find(v),o=JSON.parse(i.val()),o.Description=r.val(),this.setMediaData(a,o)},initMedia:function(o){var r=e(b),a=e(h);o&&(r=e(y).find(b),a=e(y).find(h)),e(y).find(".qor-table--medialibrary-file").each(function(){e(this).closest(".mdl-card__supporting-text").addClass("qor-table--files")}),(r.length||a.length)&&(r.each(function(){var t=e(this),i=t.data("videolink"),o=i&&i.match(/\.mp4$|\.m4p$|\.m4v$|\.m4v$|\.mov$|\.mpeg$|\.webm$|\.avi$|\.ogg$|\.ogv$/);o&&t.parent().addClass("qor-table--video qor-table--video-internal").html('<video width=100% height=100% controls><source src="'+i+'"></video>')}),a.each(function(){var o=e(this),r=o.data("videolink"),a=t(r),n=i(r);a&&o.parent().addClass("qor-table--video qor-table--video-external").html('<iframe width="100%" height="100%" src="//www.youtube.com/embed/'+a+'?rel=0" frameborder="0" allowfullscreen></iframe>'),n&&o.parent().addClass("qor-table--video qor-table--video-external").html('<iframe width=100% height=100% src="http://player.youku.com/embed/'+n+'" frameborder=0 "allowfullscreen"></iframe>')}))},setVideo:function(o){var r=e(o.target),a=r.closest("[data-tab-source]"),n=r.closest("form"),d=n.find(v),l=JSON.parse(d.val()),s=r.val(),f=a.find("iframe"),c=t(s),u=i(s);l.SelectedType="video_link",l.Video=s,this.setMediaData(n,l),(c||u)&&(f.length&&f.remove(),c&&a.append('<iframe width="100%" height="400" src="//www.youtube.com/embed/'+c+'?rel=0" frameborder="0" allowfullscreen></iframe>'),u&&a.append('<iframe width=100% height=400 src="http://player.youku.com/embed/'+u+'" frameborder=0 "allowfullscreen"></iframe>'))},resetMediaData:function(t,i,o){var r=e(i),a=r.closest("form"),n=r.find(v),d=r.find(u).find(".qor-fieldset__alert"),l=JSON.parse(n.val());l.SelectedType=o,"video_link"==o&&(l.Video=r.find(m).val(),d.length&&d.remove()),l.Description=e('[data-tab-source="'+o+'"]').find(p).val(),e(c).val(o),this.setMediaData(a,l)},destroy:function(){this.unbind()}},o.DEFAULTS={},e.fn.qorSliderAfterShow=e.fn.qorSliderAfterShow||{},e.fn.qorSliderAfterShow.renderMediaVideo=function(){var o=e(u),r=e(p),a=o.length&&o.data().videourl,n=a&&t(a),d=a&&i(a);r.length&&r.val(r.data().imageInfo.Description),o.length&&a&&(n&&o.append('<iframe width="100%" height="400" src="//www.youtube.com/embed/'+n+'?rel=0&fs=0&modestbranding=1&disablekb=1" frameborder="0" allowfullscreen></iframe>'),d&&o.append('<iframe width=100% height=400 src="http://player.youku.com/embed/'+d+'" frameborder=0 "allowfullscreen"></iframe>'))},o.plugin=function(t){return this.each(function(){var i,a=e(this),n=a.data(r);if(!n){if(/destroy/.test(t))return;a.data(r,n=new o(this,t))}"string"==typeof t&&e.isFunction(i=n[t])&&i.apply(n)})},e(function(){var t=".qor-table--medialibrary";e(document).on(n,function(i){o.plugin.call(e(t,i.target),"destroy")}).on(a,function(i){o.plugin.call(e(t,i.target))}).triggerHandler(a)}),o});