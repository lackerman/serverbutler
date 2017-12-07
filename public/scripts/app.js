function getData(url, callback) {
  const request = new XMLHttpRequest();
  request.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
		if (this.status == 200) {
          callback(request.responseText, null);
		} else {
		  callback(null, this.status);
		}
	}
  };
  request.open("GET", url, true);
  request.send();
}
