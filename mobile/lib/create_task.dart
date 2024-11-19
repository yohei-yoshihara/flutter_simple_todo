import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:fluttertoast/fluttertoast.dart';

import 'package:simple_todo/model.dart';
import 'package:simple_todo/session.dart';
import 'package:simple_todo/logging.dart';
import 'package:simple_todo/datepicker.dart';
import 'package:simple_todo/loading.dart';

class CreateTaskPage extends StatefulWidget {
  final Function(Task)? onCreated;

  const CreateTaskPage({super.key, required this.onCreated});

  @override
  CreateTaskPageState createState() => CreateTaskPageState();
}

class CreateTaskPageState extends State<CreateTaskPage> {
  final _formKey = GlobalKey<FormState>();
  final _titleController = TextEditingController();
  final _descriptionController = TextEditingController();
  bool _completed = false;
  final DateTime _dueDate = DateTime.now();
  bool _isLoading = false;

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
      "/api/tasks/create",
    );
    String body = jsonEncode({
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

    log.info('create: Response status: ${response.statusCode}');
    log.info('create: Response body: ${response.body}');

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

    if (response.statusCode != 200 && context.mounted) {
      Fluttertoast.showToast(
        msg: "タスクの作成に失敗しました",
        toastLength: Toast.LENGTH_SHORT,
        gravity: ToastGravity.BOTTOM,
        backgroundColor: Colors.red.shade300,
        textColor: Colors.white,
      );
      return;
    }

    if (widget.onCreated != null && context.mounted) {
      var json = jsonDecode(response.body);

      var task = Task(
        json["id"],
        _titleController.text,
        _descriptionController.text,
        _completed,
        _dueDate,
      );
      widget.onCreated!(task);
      Navigator.of(context).pop();
    }
  }

  @override
  void dispose() {
    _titleController.dispose();
    _descriptionController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("新しいタスクの作成")),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
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
                      onChanged: (value) {},
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
                      child: const Text("新しいタスクの作成"),
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
