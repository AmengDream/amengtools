// ==UserScript==
// @name         自动看课
// @namespace    http://tampermonkey.net/
// @version      1.0
// @description  自动点击“学习下一节”按钮并将播放器进度调到最后
// @author       ameng
// @match        https://whqgu.ls365.net/*
// @grant        none
// ==/UserScript==

(function() {
    'use strict';

    // 定义检查并点击按钮的函数
    function checkAndClick() {
        const successDiv = document.querySelector('#reader_success_video.success');
        if (successDiv && successDiv.style.display === 'block') {
            const nextButton = document.getElementById('learnNextSection');
            if (nextButton) {
                nextButton.click();
            }
        }
    }

    // 定义调整播放器进度的函数
    function adjustVideoProgress() {
        const player = document.querySelector('video');
        if (player) {
            // 将播放器进度设置到最后
            player.currentTime = player.duration;
        }
    }

    // 在页面加载时运行该函数
    window.addEventListener('load', () => {
        checkAndClick();
        adjustVideoProgress();
    });

    // 可选：设置一个定时器，每隔一段时间检查一次
    setInterval(() => {
        checkAndClick();
        adjustVideoProgress();
    }, 1000);
})();