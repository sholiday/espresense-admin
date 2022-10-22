$(document).ready(function () {
  $.getJSON("/table-dev", function (data) {
    var items = [
      { data: "mac" },
      { data: "manuf" },
      { data: "idtype" },
      { data: "name", name: "Device" },
      { data: "disc" },
      { data: "closest" },
    ];

    for (i = 0; i < data.rooms.length; i++) {
      items.push({ data: data.rooms[i] });
    }

    for (i = 0; i < items.length; i++) {
      $("#col-head").append("<th>" + items[i].data + "</th>");
      $("#col-footer").append("<th>" + items[i].data + "</th>");
    }

    var table = $("#example").DataTable({
      scrollY: "500px",
      scrollCollapse: true,
      paging: false,
      ajax: {
        url: "/table-dev",
      },
      columns: items,
    });
    setInterval(function () {
      if ($("#auto-refresh-chk").is(":checked")) {
        table.ajax.reload(null, false); // user paging is not reset on reload
      }
    }, 500);
  });
});

$(document).ready(function () {
  $.getJSON("/table-rooms", function (data) {
    var items = [
      { data: "name" },
      { data: "IP" },
      { data: "Uptime" },
      { data: "Firm" },
      { data: "Ver" },
      { data: "Rssi" },
      { data: "Adverts" },
      { data: "Seen" },
      { data: "Queried" },
      { data: "Reported" },
    ];

    for (i = 0; i < items.length; i++) {
      $("#rcol-head").append("<th>" + items[i].data + "</th>");
      $("#rcol-footer").append("<th>" + items[i].data + "</th>");
    }

    var table1 = $("#rooms").DataTable({
      scrollY: "500px",
      scrollCollapse: true,
      paging: false,
      ajax: {
        url: "/table-rooms",
      },
      columns: items,
    });
    setInterval(function () {
      if ($("#auto-refresh-chk").is(":checked")) {
        table1.ajax.reload(null, false); // user paging is not reset on reload
      }
    }, 2000);
  });
});
