function edit() {
  $.ajax({
    type: "get",
    url: "/views",
    cache: false,
    dataType: "json",
    success: function (data) {
      $("#edit_set_table").html(
        "<tr><th>id</th><th width='7%'>本地端口</th><th width='7%'>心跳端口</th>" +
        "<th width='12%'>本地ip</th><th width='12%'>远端ip</th>" +
        "<th width='12%'>抓包时间(min)</th><th width='14%'>文件限制大小(MB)</th>" +
        "<th width='11%'>文件存储路径</th><th width='8%'>文件名</th>" +
        "<th width='14%'>操作</th></tr>"
      );
      $.each(data, function (index, item) {
        $("#edit_set_table").append(
          "<tr><td class='set_id'></td>" +
          "<td class='nativeServerPort'></td>" +
          "<td class='heart_port'></td>" +
          "<td class='nativeIp'></td>" +
          "<td class='remoteIp'></td>" +
          "<td class='packetHoldingTime'></td>" +
          "<td class='fileMaxSize'></td>" +
          "<td class='savePath'></td>" +
          "<td class='fileName'></td>" +
          "<td><input type='button' value='修改' id='modify_button' class='btn btn-success' data-toggle='modal'>&nbsp;&nbsp;<input type='button' value='删除' id='delete_button' class='btn btn-danger'></td></tr>");
        $("#edit_set_table .set_id")[index].innerText = item.setId;
        $("#edit_set_table .nativeServerPort")[index].innerText = item.nativeServerPort;
        $("#edit_set_table .heart_port")[index].innerText = item.heart_port;
        $("#edit_set_table .nativeIp")[index].innerText = item.nativeIp;
        $("#edit_set_table .remoteIp")[index].innerText = item.remoteIp;
        $("#edit_set_table .packetHoldingTime")[index].innerText = item.packetHoldingTime;
        $("#edit_set_table .fileMaxSize")[index].innerText = item.fileMaxSize;
        $("#edit_set_table .savePath")[index].innerText = item.savePath;
        $("#edit_set_table .fileName")[index].innerText = item.fileName;
      });
      $("#edit_set").modal("show");
    },
    error: function () {
      alert("无法连接到服务器");
    }
  });

}


