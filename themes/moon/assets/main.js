$.ajaxSetup({
  beforeSend: function(xhr) {
    xhr.setRequestHeader('X-Xsrftoken', $('meta[name=_xsrf]').attr('content'));
  }
});

$(function() {
  $("form[data-confirm]").submit(function(e) {
    if(!confirm($(this).data('confirm'))){
      e.preventDefault();
    }
  });

  $("a[data-method]").click(function(e) {
    e.preventDefault();
    var msg = $(this).data('confirm');
    var method = $(this).data('method');
    var next = $(this).data('next');
    var url = $(this).attr('href');

    var ok = true;
    if (msg) {
      if (!confirm(msg)) {
        ok = false;
      }
    }
    if (ok) {
      // console.log(method, url, next);
      $.ajax({
        type: method,
        url: url,
        success: function() {
          window.location.href = next;
        }
      })
    }
  });
});