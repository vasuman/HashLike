function supportsEmscripten() {
    //FIXME!!
    return true;
}

function fuzzyRepr(num) {
    // FIXME!!
    return "" + num;
}

/* Required options
   - url
   - sys
   - countTarget
   - serverAddr
*/

function HashLike(options) {
    this.opt = options;
    this.getLikes();
}

HashLike.prototype.getLikes = function () {
    var self = this,
        urlParam = 'url=' + encodeURIComponent(this.opt.url),
        powParam = 'sys=' + this.opt.sys,
        params = '?' + urlParam + '&' + powParam,
        countURL = this.opt.serverAddr + '/count' + params;
    function loadCback() {
        if (this.status != 200) {
            console.log('Error!!', this.responseText);
            return;
        }
        var count;
        try {
            count = Integer.parseInt(this.responseText);
        } catch (e) {
            console.log('Not a number', this.responseText);
            return;
        }
        var targetElt = countTarget;
        targetElt.innerHTML = fuzzyRepr(this.responseText);
        // TODO: Set tooltip text
    }
    var xhr = new XMLHttpRequest();
    xhr.open('GET', countURL, true);
    xhr.onload = loadCback;
    xhr.send();
}

HashLike.prototype.initWorker = function () {
    if (!window.Worker) {
        throw new Error("Browser doesn't support Web Workers!");
    }
    var workerFile;
    switch (this.opt.powSystem) {
    case "HC":
        workerFile = 'hc-worker';
        break;
    default:
        throw new Error('Unsupported proof-of-work system');
    }
    if (supportsEmscripten()) {
        workerFile += '-emc';
    }
    this.worker = new Worker(workerFile + '.js');
    this.worker.postMessage('hello');
}
