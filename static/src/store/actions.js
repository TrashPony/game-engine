import {urls} from "../const";

const actions = {
  sendSocketData({commit, getters}, payload) {

    let tryCount = 5;

    let waitConnect = setInterval(function () {
      if (getters.getWSConnectState.connect) {
        getters.getWSConnectState.ws.send(payload);
        clearInterval(waitConnect)
      } else {
        if (!getters.getWSConnectState.pending) {
          commit({type: 'reconnectWS'});
        }
      }

      tryCount--
      if (tryCount === 0) {
        clearInterval(waitConnect)
      }
    }, 100)
  },
  changeSettings({commit, getters}, payload) {
    getters.getWSConnectState.ws.send(JSON.stringify({
      event: "changeSettings",
      name: payload.name,
      service: "system",
      count: Number(Math.round(payload.count))
    }));
  },
  getAvatar(context, payload) {

    let avatars = context.getters.getChatState.avatars;
    if (avatars.hasOwnProperty(payload.userID)) return;

    context.commit({
      type: 'addAvatar',
      id: payload.userID,
      avatar: "url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGAAAABgCAYAAADimHc4AAAABmJLR0QA/wD/AP+gvaeTAAAIi0lEQVR4nO2cbXBU1RnH/8+9mxdqojhtGBGZFpOYhARIQhypEmOgyEtTgjDU9kMpkQQcW3FaoigNdAspkqKtnamjhCg6jtihlYgpxmGsG4xvU9a8VEk2zUYQK1ojRCAQkt29Tz+QfEjc7N5zd0+ymzm/j/c8zzn/uf+995x79pwDKBQKhUKhUCgUCoVCoVCMGTTeAkJh51tLk2KBAgPaVIAv6Mytmwrqm8dblwhRacAf3y2a5vVyFRH/FIA2orhdI3p40/zDr46HNlGizoDHG4vmGDDqAUwNEMYAtj+Y/5p9jGRZJqoM2PnW0qQYoiYAN5iJJ6K15fMPPy9ZVkiMfHwjmhiiCpi8+QDAzLur3l6eKFFSyESNAXucc2MArBFMS9LZu0KGnnARNQb0XpySBWCyaJ4Bni9BTtiIGgNYp+9YySNYyxsrbOMtwCw+oMfKr4UJZ0Xi7QdWxxpxZwsMDfkg43oYdImJ3aTT4cplb3ZZkBCQqDFg0nn+sD9R6wU4QSSPQO+bjd1aV7jGi68qAZoOAGACCCAQ4MMTW+sWvKKTVm4veuNjQfkB9I0R+5zrZxnMNxlgr82mt5Xk7OkUrWN347IaAOsEUs57YuNu3HJL7ZlAQXa7XfPkNj5FxOtN1HmWmFfuWO44KqBjVKQbsPdY6SIi+hOAzGEFzO+Rrm9cl7vHabaux4/eOd0gvRlE3zaVQPyrB+fXPxEsrKJu4U6AHzGrA8DXMHheZbGjQyDHL1I74ZpjZfcQUT1G3nwAIPo+G0bjM8fKlpqtb1PBkU8ZWAXgfLBYBleX31b/52BxWw4tnA3wZrMaBpkMjfYK5vhFmgE1x9ZlgfA0AD1AWDwTXtr3r/uuM1vvQ7fXHyVdvwXAG6OEfA6i0ofy6zcQgYPVp+n8G1i7D/lbDxXeZiFvGPI6YaJHAMSYiLzG0AYeAGD6FVB+a50LwKKqoz/MII0XEHgqE/VqBjXH9hoNG5e91m+mHvuB1bFePmP6CfwGGpYDeMdyPuSOghabDWTQEggYMMTmgsPtANpF84bwJpyZDh8sT1UwaKbV3CGkvIL2OdbGAyY7yitMk6EjGAYboc4TXROqBikGlBQ+d5mBPoGUHhk6gkGG1h1SPujLUDXIGwWRwLuR8LY0HQH4/Y/e/AzAJ1bzmfBeqBqkGcBA0CHgIIbG9KQsHcGhAxYTvTD0g6G2Ls2A9XP3/gNA0LEyg+z35FU3ydIRDJ8+sBvAOdE8Bmoqlx85EWr7Uj/E/jt32r0E3gHA37DwApjuL8ur3iFTQzAeXdbYrRHfDcAnkNbaH2srD0f7YzIXtPf90hs0G93FhJuI2QDRR16m2g151V+NRftmqHh1YRGIXwRwdcBAwjuGj1btLP7n/8LRblT9Jywbe+0d3/PEaL8lxmoAV40o7iTix/Sr+Fl7YYM3XG0qA/xgd9wR77mo5RL4OgPoj/Vxh724wT3euhQKhUKhUCgUijAR8nfAy22V39XgWwFGMoM9DO0Dmx5XV5y++UI4BE50LBvgcNhtPVNQBdD9+OZfj90g2rhy5ra/hiZv4mNpMo6ZqWcK7Qfo1/D/v28SmPe/3La9NDR5Ex9LT0Dt8e1rGbzPROhlLzjtx5n2U1baiXSaXZ3FDMojhtt38ev9eXl5HtE6rD0B4AdMhsbroHuttBHptLi6agB6hYAKEJ7TEyY3OJ1OM6tAhiFswCFXVSKAbLPxBOSLthHpNHd03czgkUskb7UlXPsz0bqEDWCPR3SN/rWibUQ6xEaav+sMThetS9iAngRvNwCR+fDTom1EPD40AzBGXiaw6XWuQwgbUDLDfplADeYz6HXRNiKd7MzU4wzYMdyEA3PSUv4mWpe1TlhDJRB83SWAz3U9LiyLWCON3PSUHTpTNhglRJSfk55yNxGZuSfDsPwhdvD49ocBfjRASC8zLVmVtS2ktZMTHcurIlZmbtvF4NUATvopdhigeermByfkuSA727Wsj5CraZQGg/uY0bJqtj1sW3gUCoVCoVAoFAqFQqEIK2OyOrrV5V5gMH7ChHQCDAJaDdALuenJwtO3Ew2pBrhcrsQ+tj0Pwl1+ipmB6jhf/8bMzMwBmToiGWkbtR0Oh62PbQdB+MEoIUTAhgEtbhKAn8vSEelIewKaOzpLweYOtCDmxdkZqUdkaRHB6Tz9LT3h0i9AfDszxYDo3Thv7F8yM6cLHfxkFnmb9BhlpkNJi4j1Qy6XK1FLuNQI4A9gKiJgMTH/bkAfcDrbTgY6p9QyUgxgZg2gXIGMm2XoEKUPtgoC/OjmGbrmeUxGm1IMcLvdMRDoXwiIl6FDGMai0QvpThlNSjEgNTW1H8BnZuMZCPtheJaggEtoEh0OR9gHLdL6ACKY3sbPjFpZOkQg8AejFjK1FBYWhm176hASd8r7qgA2MXKgE/rAxafl6TCPAW0X/O/qZ+iQsqNf6odYU3tnAREdwujn6pzWmZbMzkj+0Er9rW3uLEPDGgBpBJwDyBHju/xiKB92LS53IQNPAsgYvPQJg8tz01P/brXOQEifinB2dibrBu0AYwWASYOXzxHwku7R7LNm3Whpy3+Ty72VriyOGv4UM9o08hbNSU8P6SCNps7OJM0wYrPT0kz3ZVYYs53yjhMn4q/2eGboPt3X88Wpj0N5n7a4ukoY/GyAkPYEnXMGBwMRTdQdVcDMWktH1ykEO+aMuCwnLbVmbFRZJ2oO7x7i3+1dM2HijDlmGm0OKqKIOgMMCnKczCAU7NiZCCHqDNBi4IafpeF++I9sLeEg6gyYk5LyJRj1QcLYYLwwJoJCJOoMuIL+SwCjHxlJtGtuRsroX7URRFQakJMx46QG7zwCXsfwfQrdDLovJy15y3hpEyXqhqEjaXK5rgfrqQzqPf/Fp60y5msUCoVCoVAoFAqFQqFQhIv/A1oSm/LE8/ruAAAAAElFTkSuQmCC')",
    });

    let time = new Date().getTime()
    payload.vm.$api.get(urls.avatarURL + '?user_id=' + payload.userID + "&timestamp=" + time + "").then(function (response) {

      let loadAvatar = function (base64) {
        context.commit({
          type: 'addAvatar',
          id: payload.userID,
          type_avatar: response.data.type,
          avatar: "url('" + base64 + "')",
        });
      }

      if (response.data.type === "base64") {
        loadAvatar(response.data.avatar)
      }

      if (response.data.type === "path") {
        toDataURL(require('../assets/' + response.data.avatar), loadAvatar)
      }
    });
  },
  getNpcAvatar(context, payload) {
    let avatars = context.getters.getChatState.avatars;
    let id = 'dialog:' + payload.page_id + payload.user_id;
    if (avatars.hasOwnProperty(id)) return;

    context.commit({
      type: 'addAvatar',
      id: id,
      avatar: "url('data:image/png;base64,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')",
    });

    payload.vm.$api.get(urls.dialogPicURL + '?dialog_page_id=' + payload.page_id + "&user_id=" + payload.user_id).then(function (response) {
      context.commit({
        type: 'addAvatar',
        id: id,
        type_avatar: response.data.type,
        avatar: "url('" + response.data.picture + "')",
      });
    });
  },
  playSound({getters}, payload) {

    if (payload.k === 0 || !payload.k) payload.k = 1
    if (!payload.type || payload.type === '') payload.type = 'interface';

    let path = require('assets/audio/sound_effects/' + payload.type + '/' + payload.sound)
    let audio = new Audio(path);
    audio.volume = getters.getSettings.SFXVolume * payload.k;
    audio.play();
  }
};

function toDataURL(url, callback) {
  let xhr = new XMLHttpRequest();
  xhr.onload = function () {
    let reader = new FileReader();
    reader.onloadend = function () {
      callback(reader.result);
    }
    reader.readAsDataURL(xhr.response);
  };
  xhr.open('GET', url);
  xhr.responseType = 'blob';
  xhr.send();
}


export default actions;
