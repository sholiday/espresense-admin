$(document).ready(function() {
  $.getJSON("/table", function(data) {
    var items = [ {data : 'name', name : 'Device'}, {data : 'closest'} ];
    for (i = 0; i < data.rooms.length; i++) {
      items.push({data : data.rooms[i]});
    }

    for (i = 0; i < items.length; i++) {
      $("#col-head").append('<th>' + items[i].data + '</th>');
      $("#col-footer").append('<th>' + items[i].data + '</th>');
    }

    var table = $('#example').DataTable({
      scrollY : '500px',
      scrollCollapse : true,
      paging : false,
      ajax : {
        "url" : "/table",
      },
      columns : items,
    });
    setInterval(function() {
      table.ajax.reload(null, false); // user paging is not reset on reload
    }, 500);
  });
});