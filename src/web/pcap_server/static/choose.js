$(function () {
  $("#start_button").click(function () {
    var flag = false;
    var target = "/start";
    var checks = document.getElementsByName("choose");
    for (var i = 0; i < checks.length; i++) {
      if (checks[i].checked == true) {
        var id = $(".choose").eq(i).parent().parent().find(".set_id").html().trim();
        if (!flag) {
          target = target + "?id=" + id;
        } else {
          target = target + "&id=" + id;
        }
        flag = true;
      }
    }
    if (!flag) {
      alert("请至少选择一项");
    } else {
      $("#choose_set").modal("hide");
      $.ajax({
        url: target,
        method: "get",
        dataType: "json",
        success: function (data) {
          $("#state").html(
            "<tr><th width='8%'>设置id</th><th width='14%'>已抓包数量</th><th width='20%'>已抓文件大小(MB)</th><th width='14%'>已存文件数量</th><th width='15%'>已抓包时间</th><th width='17%'>抓包状态</th><th width='12%'>操作</th></tr>"
          );
          $.each(data.Message, function (i, item) {
            if (item.Finished) {
              $("#state").append(
                "<tr><td class='set_id' data-toggle='modal' data-target='#view_set_id'> " +
                item.SetId +
                "</td>" +
                "<td class='packet_count'> " +
                item.PacketCount +
                "</td>" +
                "<td class='file_size'> " +
                item.FileSize.toFixed(3) +
                "</td>" +
                "<td class='file_count'> " +
                item.FileCount +
                "</td>" +
                "<td class='hold_time'> " +
                formatSeconds(item.PGHoldedTime) +
                "</td>" +
                "<td class='finish'>  " +
                item.State +
                "</td><td><input type='button' value='停止抓包' disabled='true' id='end_button' class='btn btn-danger'></td></tr>"
              );
            } else {
              $("#state").append(
                "<tr><td class='set_id' data-toggle='modal' data-target='#view_set_id'> " +
                item.SetId +
                "</td>" +
                "<td class='packet_count'> " +
                item.PacketCount +
                "</td>" +
                "<td class='file_size'> " +
                item.FileSize.toFixed(3) +
                "</td>" +
                "<td class='file_count'> " +
                item.FileCount +
                "</td>" +
                "<td class='hold_time'> " +
                formatSeconds(item.PGHoldedTime) +
                "</td>" +
                "<td class='finish'> " +
                item.State +
                "</td><td><input type='button' value='停止抓包'  id='end_button' class='btn btn-danger'></td></tr>"
              );
            }
          });
        }
      });
    }
  });
});
function start() {
  $("#choose_table").html(
    "<tr><th width='4%'>id</th><th width='7%'>本地端口</th><th width='7%'>心跳端口</th>" +
    "<th width='14%'>本地ip</th><th width='14%'>远端ip</th>" +
    "<th width='12%'>抓包时间(min)</th><th width='15%'>文件限制大小(MB)</th>" +
    "<th width='12%'>文件存储路径</th><th width='10%'>文件名</th>" +
    "<th width='5%'>选择</th></tr>"
  );
  $.ajax({
    url: "/rule",
    method: "get",
    dataType: "json",
    success: function (data) {
      if (data == null) {
        $("#start_button").attr("disabled", true);
      } else {
        $("#start_button").attr("disabled", false);
      }
      $.each(data, function (i, item) {
        var set = item.Set;
        if (item.Flag) {
          $("#choose_table").append(
            "<tr class='special'>" +
            "<td class='set_id'></td>" +
            "<td class='nativeServerPort'></td>" +
            "<td class='heart_port'></td>" +
            "<td class='nativeIp'></td>" +
            "<td class='remoteIp'></td>" +
            "<td class='packetHoldingTime'></td>" +
            "<td class='fileMaxSize'></td>" +
            "<td class='savePath'></td>" +
            "<td class='fileName'></td>" +
            "<td><input type='checkbox' name='choose' class='choose'  disabled='true'/></td></tr>"
          );
        } else {
          $("#choose_table").append(
            "<tr><td class='set_id'></td>" +
            "<td class='nativeServerPort'></td>" +
            "<td class='heart_port'></td>" +
            "<td class='nativeIp'></td>" +
            "<td class='remoteIp'></td>" +
            "<td class='packetHoldingTime'></td>" +
            "<td class='fileMaxSize'></td>" +
            "<td class='savePath'></td>" +
            "<td class='fileName'></td>" +
            "<td><input type='checkbox' name='choose' class='choose'/></td></tr>"
          );
        }
        $("#choose_table .set_id")[i].innerText = set.setId;
        $("#choose_table .heart_port")[i].innerText = set.heart_port;
        $("#choose_table .nativeServerPort")[i].innerText = set.nativeServerPort;
        $("#choose_table .nativeIp")[i].innerText = set.nativeIp;
        $("#choose_table .remoteIp")[i].innerText = set.remoteIp;
        $("#choose_table .packetHoldingTime")[i].innerText = set.packetHoldingTime;
        $("#choose_table .fileMaxSize")[i].innerText = set.fileMaxSize;
        $("#choose_table .savePath")[i].innerText = set.savePath;
        $("#choose_table .fileName")[i].innerText = set.fileName;
      });
      $("#choose_set").modal("show");
    },
    error: function () {
      alert("无法连接到服务器");
    }
  });
}
