let msgCalc = {};
let times = [];

setInterval(function () {
  if (Object.keys(msgCalc).length !== 0) {

    let allCount = 0;
    let allSize = 0;

    for (let i in msgCalc) {

      if (Math.round(msgCalc[i].size / 1024) >= 1) {
        console.log(i, "count: ", msgCalc[i].count, ", size: ", Math.round(msgCalc[i].size / 1024), 'kb');
      }

      allCount += msgCalc[i].count;
      allSize += msgCalc[i].size;
    }

    const sum = times.reduce((a, b) => a + b, 0);
    const avg = (sum / times.length) || 0;

    console.log("all - count: ", allCount, ", size: ", Math.round(allSize / 1024), 'kb, avg: ', avg);
    console.log("-----------");
    console.log();


    times.push(Math.round(allSize / 1024))

    msgCalc = {}
  }
}, 1000);

function logMsg(event, data) {
  const blob = new Blob([JSON.stringify(data)], {type: 'application/json'});
  if (!msgCalc.hasOwnProperty(event)) msgCalc[event] = {count: 0, size: 0};
  msgCalc[event].count++;
  msgCalc[event].size += blob.size;
}

function logBinMsg(event, blob) {
  if (!msgCalc.hasOwnProperty(event)) msgCalc[event] = {count: 0, size: 0};
  msgCalc[event].count++;
  msgCalc[event].size += blob.size;
}

function logArrayMsgEvent(event, array) {
  if (!msgCalc.hasOwnProperty(event)) msgCalc[event] = {count: 0, size: 0};
  msgCalc[event].count++;
  msgCalc[event].size += array.length;
}

export {logMsg, logBinMsg, logArrayMsgEvent}
