import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:fluttertoast/fluttertoast.dart';
import 'package:simple_todo/loading.dart';

import 'package:simple_todo/model.dart';
import 'package:simple_todo/session.dart';
import 'package:simple_todo/logging.dart';
import 'package:simple_todo/datepicker.dart';

class EditTaskPage extends StatefulWidget {
  final Task task;
  final Function(Task)? onUpdated;
  final Function(Task)? onDeleted;
  const EditTaskPage({super.key, required this.task, required this.onUpdated, required this.onDeleted});

  @override
  State<EditTaskPage> createState() => _EditTaskPageState();
}

class _EditTaskPageState extends State<EditTaskPage> {
  final _formKey = GlobalKey<FormState>();
  final _titleController = TextEditingController();
  final _descriptionController = TextEditingController();
  bool _completed = false;
  DateTime _dueDate = DateTime.now();
  bool _isLoading = false;

  @override
  void initState() {
    _titleController.text = widget.task.title;
    _descriptionController.text = widget.task.description;
    _completed = widget.task.completed;
    _dueDate = widget.task.dueDate;

    super.initState();
  }

  @override
  void dispose() {
    _titleController.dispose();
    _descriptionController.dispose();
    super.dispose();
  }

  void _notifyTaskChanged() {
    if (widget.onUpdated != null) {
      var task = Task(
        widget.task.id,
        _titleController.text,
        _descriptionController.text,
        _completed,
        _dueDate,
      );
      widget.onUpdated!(task);
    }
  }

  void _submit(BuildContext context) async {
    if (!context.mounted) return;

    if (_formKey.currentState != null && !_formKey.currentState!.validate()) {
      return;
    }

    setState(() {
      _isLoading = true;
    });

    var url = Uri.http(
      "localhost:8000",
      "/api/tasks/update",
    );
    String body = jsonEncode({
      "id": widget.task.id,
      "title": _titleController.text,
      "description": _descriptionController.text,
      "completed": _completed,
      "due_date": _dueDate.toUtc().toIso8601String(),
      "folder_id": 1
    });

    var store = SessionStore();
    http.Response response = await http.post(url, body: body, headers: store.headers);
    store.updateCookie(response);

    setState(() {
      _isLoading = false;
    });

    log.info('update: Response status: ${response.statusCode}');
    log.info('update: Response body: ${response.body}');

    if (response.statusCode == 401 && context.mounted) {
      Fluttertoast.showToast(
        msg: "セッションがタイムアウトしました。再度ログインしてください。",
        toastLength: Toast.LENGTH_LONG,
        gravity: ToastGravity.BOTTOM,
        backgroundColor: Colors.red.shade300,
        textColor: Colors.white,
      );
      Navigator.of(context).popUntil((route) => route.isFirst);
      Navigator.pushReplacementNamed(context, "/login");
      return;
    }

    if (response.statusCode != 200) {
      Fluttertoast.showToast(
        msg: "タスクの更新に失敗しました",
        toastLength: Toast.LENGTH_SHORT,
        gravity: ToastGravity.BOTTOM,
        backgroundColor: Colors.red.shade300,
        textColor: Colors.white,
      );
      return;
    }
    _notifyTaskChanged();
    if (context.mounted) {
      Navigator.of(context).pop();
    }
  }

  void _onDelete(BuildContext context) async {
    var url = Uri.http(
      "localhost:8000",
      "/api/tasks/delete",
    );
    String body = jsonEncode({
      "id": widget.task.id,
    });

    var store = SessionStore();
    http.Response response = await http.post(url, body: body, headers: store.headers);
    store.updateCookie(response);

    log.info('Response status: ${response.statusCode}');
    log.info('Response body: ${response.body}');

    if (response.statusCode != 200 && context.mounted) {
      Fluttertoast.showToast(
        msg: "タスクの削除に失敗しました",
        toastLength: Toast.LENGTH_SHORT,
        gravity: ToastGravity.BOTTOM,
        backgroundColor: Colors.red.shade300,
        textColor: Colors.white,
      );
      return;
    }

    if (context.mounted) {
      Fluttertoast.showToast(
        msg: "タスクを削除しました",
        toastLength: Toast.LENGTH_SHORT,
        gravity: ToastGravity.BOTTOM,
        backgroundColor: Colors.green.shade300,
        textColor: Colors.white,
      );
      if (widget.onDeleted != null) {
        widget.onDeleted!(widget.task);
      }
      Navigator.of(context).pop();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("タスクの詳細"),
        actions: [
          PopupMenuButton<String>(
            itemBuilder: (context) => [
              PopupMenuItem(
                value: "Delete",
                child: const Text("削除"),
                onTap: () {
                  _onDelete(context);
                },
              ),
            ],
          ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: PopScope(
          canPop: true,
          child: Stack(
            children: [
              Form(
                key: _formKey,
                child: Column(
                  children: [
                    TextFormField(
                      controller: _titleController,
                      decoration: const InputDecoration(labelText: "タイトル"),
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return "タイトルを入力してください";
                        }
                        return null;
                      },
                    ),
                    TextFormField(
                      keyboardType: TextInputType.multiline,
                      maxLines: null,
                      controller: _descriptionController,
                      decoration: const InputDecoration(labelText: "説明"),
                    ),
                    DatePicker(
                      dueDate: _dueDate,
                      onChanged: (value) {
                        setState(() {
                          _dueDate = value;
                        });
                      },
                    ),
                    SwitchListTile(
                      title: const Text("完了?"),
                      value: _completed,
                      onChanged: (value) {
                        setState(() {
                          _completed = value;
                        });
                      },
                    ),
                    ElevatedButton(
                      onPressed: () {
                        _submit(context);
                      },
                      child: const Text("タスクの更新"),
                    )
                  ],
                ),
              ),
              LoadingWidget(isLoading: _isLoading)
            ],
          ),
        ),
      ),
    );
  }
}
