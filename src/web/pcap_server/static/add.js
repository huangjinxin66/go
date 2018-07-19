  $(function () {
            $("#add_set_table #save_button").click(function () {
			    if ($("#add_set_table .port").val() == "") {
                    alert("本地端口号不能为空");
                    return;
                }
                if ($("#add_set_table .port").val() != "" && !isLegalPort($("#add_set_table .port").val())) {
                    alert("请输入合法的本地端口号：0-65535");
                    return;
                }
                if ($("#add_set_table .heart_port").val() == "") {
                    alert("心跳端口不能为空");
                    return;
                }
                if ($("#add_set_table .heart_port").val() != "" && !isLegalPort($("#add_set_table .heart_port").val())) {
                    alert("请输入合法的心跳端口号：0-65535");
                    return;
                }
                if ($("#add_set_table .heart_port").val() ==$("#add_set_table .port").val()) {
                    alert("心跳端口号和本地端口号不能重复");
                    return;
                }
                if ($("#add_set_table .during_time").val() != ""&&!isNum($("#add_set_table .during_time").val())) {
                    alert("抓包时间请输入正整数");
                    return;
                }
                if (!$("#add_set_table .file_max_size").val() == "" && !(isDecimal($("#add_set_table .file_max_size").val()) || isNum($("#add_set_table .file_max_size").val()))) {
                    alert("文件限制大小请输入正整数或小数");
                    return;
                }
                if ($("#add_set_table .file_name").val() == "") {
                    alert("请输入文件名");
                    return;
                }
                if (!checkFileName($("#add_set_table .file_name").val())) {
                    alert("文件名不能有斜杠");
                    return;
                }
                var data = $("#add_set_form").serialize();
                $.ajax({
                    type: 'post',
                    url: '/set',
                    cache: false,
                    data: data,
                    dataType: 'text',
                    beforeSend: function () {
                       $("#load").show();
                    },  
                    success: function (data) { 
                       $("#load").hide();
                        alert(data);
                        if (data == "保存配置成功") {
                            $("#add_set_form")[0].reset()
                            $("#add_set").modal('hide');
                        }
                    },
                    error: function () {
                       $("#load").hide();
                        alert("无法连接到服务器")
                    }
                });
            });
        });