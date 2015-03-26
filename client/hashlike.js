/* Required options
   - url
   - powSystem
   - elt
   - serverAddr
*/

function HashLike(options) {
    this.opt = options;
    this.getLikes();
}

HashLike.prototype.getLikes = function () {
    var self = this,
        urlParam = 'url=' + encodeURIComponent(this.opt.url),
        powParam = 'pow=' + this.opt.powSystem,
        params = '?' + urlParam + '&' + powParam,
        countURL = this.opt.serverAddr + '/count' + params;
    function loadCback() {
        if (this.status != 200) {
            console.log('Error!!', this.responseText);
            return;
        }
        // TODO: Fix -- insecure!!
        self.opt.elt.innerHTML = this.responseText;
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
        workerFile = 'hashcash-worker';
        break;
    default:
        throw new Error('Unsupported proof-of-work system');
    }
    this.worker = new Worker(workerFile + '.js');
    this.worker.postMessage('hello');
}
