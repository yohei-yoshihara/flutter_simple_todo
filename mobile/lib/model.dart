import 'package:intl/intl.dart';

class Task {
  int id;
  String title;
  String description;
  bool completed;
  DateTime dueDate;

  Task(
    this.id,
    this.title,
    this.description,
    this.completed,
    this.dueDate,
  );

  Task.fromJson(Map<String, dynamic> json)
      : id = json["id"] as int,
        title = json["title"] as String,
        description = json["description"] as String,
        completed = json["completed"] as bool,
        dueDate = DateTime.parse(json["due_date"] as String).toLocal();
}

final dateFormat = DateFormat('yyyy-MM-dd');
