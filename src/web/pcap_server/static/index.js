$(function () {
  $.ajax({
    url: "/state",
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
  var wsurl;
  var sock;
  var heartCheck = {
    timeout: 10000,//10s
    timeoutObj: null,
    reset: function () {
      clearTimeout(this.timeoutObj);
      this.start();
    },
    start: function () {
      this.timeoutObj = setTimeout(function () {
        reconnect()
      }, this.timeout)
    }
  }
  $.ajax({
    type: "get",
    url: "/ip",
    cache: false,
    async: false,
    dataType: "text",
    success: function (data) {
      wsurl = "ws://" + data + "/result";
    },
    error: function () {
      alert("无法连接到服务器");
    }
  });
  function reconnect() {
    sock = new WebSocket(wsurl);
    sock.onopen = function () {
      heartCheck.start();
    };
    sock.onmessage = function (e) {
      heartCheck.reset();
      var data = JSON.parse(e.data);
      $("#state").html(
        "<tr><th width='8%'>设置id</th><th width='14%'>已抓包数量</th><th width='20%'>已抓文件大小(MB)</th><th width='14%'>已存文件数量</th><th width='15%'>已抓包时间</th><th width='17'>抓包状态</th><th width='12%'>操作</th></tr>"
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
    };
    sock.onclose = function () {
      location.reload();
    };
  }
  $("body").on("click", "#delete_button", function () {
    var con;
    con = confirm("你确定要删除吗?"); //在页面上弹出对话框
    if (con == true) {
      var id = $(this).parent().parent().find(".set_id").html().trim();
      $.ajax({
        type: "get",
        url: "/delete?id=" + id,
        cache: false,
        dataType: "text",
        success: function (data) {
          alert(data);
          edit();
        },
        error: function () {
          alert("无法连接到服务器");
        }
      });
    }
  });
  $("body").on("click", "#modify_button", function () {
    var id = $(this).parent().parent().find(".set_id").html().trim();
    $("#edit_set").modal("hide");
    $.ajax({
      type: "get",
      url: "/view?id=" + id,
      cache: false,
      dataType: "json",
      async: "true",
      success: function (data) {
        $("#modify_set .heart_port").val(data.heart_port);
        $("#modify_set .port").val(data.nativeServerPort);
        $("#modify_set .nativeIp").val(data.nativeIp);
        $("#modify_set .remoteIp").val(data.remoteIp);
        $("#modify_set .during_time").val(data.packetHoldingTime);
        $("#modify_set .file_max_size").val(data.fileMaxSize);
        $("#modify_set .folder").val(data.savePath);
        $("#modify_set .file_name").val(data.fileName);
        $("#modify_set .set_id").val(id);
        $("#modify_set").modal("show");
      },
      error: function () {
        alert("无法连接到服务器");
      }
    });
  });
  $("body").on("click", "#state .set_id", function () {
    var choose_id = $(this)
      .html()
      .trim();
    $.ajax({
      type: "get",
      url: "/view?id=" + choose_id,
      cache: false,
      dataType: "json",
      success: function (data) {
        $("#view_set_state_id .heart_port")[0].innerText = data.heart_port;
        $("#view_set_state_id .port")[0].innerText = data.nativeServerPort;
        $("#view_set_state_id .nativeIp")[0].innerText = data.nativeIp;
        $("#view_set_state_id .remoteIp")[0].innerText = data.remoteIp;
        $("#view_set_state_id .during_time")[0].innerText = data.packetHoldingTime;
        $("#view_set_state_id .file_max_size")[0].innerText = data.fileMaxSize + "MB";
        $("#view_set_state_id .folder")[0].innerText = data.savePath;
        $("#view_set_state_id .file_name")[0].innerText = data.fileName;
      },
      error: function () {
        alert("无法连接到服务器");
      }
    });
  });
  $("body").on("click", "#state #end_button", function () {
    var id = $(this).parent().parent().find(".set_id").html().trim();
    var button = $(this);
    $.ajax({
      type: "get",
      url: "/stop?id=" + id,
      cache: false,
      dataType: "text",
      success: function (data) {
        button.attr("disabled", true);
      },
      error: function () {
        alert("无法连接到服务器");
      }
    });
  });
  var sock = new WebSocket(wsurl);
  console.log("onload");
  sock.onopen = function () {
    heartCheck.start();
  };
  sock.onmessage = function (e) {
    heartCheck.reset();
    var data = JSON.parse(e.data);
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
  };
  sock.onclose = function () {
    location.reload();
  };
  $("#update_button").click(function () {
    if (
      $("#modify_set .port").val() != "" &&
      !isLegalPort($("#modify_set .port").val())
    ) {
      alert("请输入合法的本地端口号：0-65535");
      return;
    }
    if (
      $("#modify_set .heart_port").val() != "" &&
      !isLegalPort($("#modify_set .heart_port").val())
    ) {
      alert("请输入合法的心跳端口号：0-65535");
      return;
    }
    if (
      $("#modify_set .port").val() != "" &&
      $("#modify_set .heart_port").val() != "" &&
      $("#modify_set .port").val() == $("#modify_set .heart_port").val()
    ) {
      alert("本地端口号和心跳端口号不能一样");
      return;
    }
    if (
      $("#modify_set .during_time").val() != "" &&
      !isNum($("#modify_set .during_time").val())
    ) {
      alert("抓包时间必须是正整数");
      return;
    }
    if (
      !$("#modify_set .file_max_size").val() == "" &&
      !(
        isDecimal($("#modify_set .file_max_size").val()) ||
        isNum($("#modify_set .file_max_size").val())
      )
    ) {
      alert("文件限制大小请输入正整数或小数");
      return;
    }
    if ($("#modify_set .file_name").val() == "") {
      alert("请输入文件名");
      return;
    }
    if (!checkFileName($("#modify_set .file_name").val())) {
      alert("文件名不能有斜杠");
      return;
    }
    var data = $("#modify_set_form").serialize();
    $.ajax({
      type: "post",
      url: "/update",
      cache: false,
      data: data,
      dataType: "text",
      beforeSend: function () {
        $("#load1").show();
      },
      success: function (data) {
        $("#load1").hide();
        alert(data);
        if (data == "更新配置成功") {
          $("#modify_set").modal("hide");
          $("#modify_set_form")[0].reset();
        }
      },
      error: function () {
        $("#load1").hide();
        alert("更新配置失败");
      }
    });
  });
});
function refresh() {
  $.ajax({
    url: "/refresh",
    method: "get",
    dataType: "json",
    success: function (data) {
		$("#state").html("<tr><th width='8%'>设置id</th><th width='14%'>已抓包数量</th><th width='20%'>已抓文件大小(MB)</th><th width='14%'>已存文件数量</th><th width='15%'>已抓包时间</th><th width='17%'>抓包状态</th><th width='12%'>操作</th></tr>"
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
