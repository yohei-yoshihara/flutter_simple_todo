import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

import 'package:simple_todo/logging.dart';

final dateFormat = DateFormat('yyyy-MM-dd');

class DatePicker extends StatefulWidget {
  final DateTime dueDate;
  final void Function(DateTime) onChanged;

  const DatePicker({
    super.key,
    required this.dueDate,
    required this.onChanged,
  });

  @override
  State<DatePicker> createState() => _DatePickerState();
}

class _DatePickerState extends State<DatePicker> {
  DateTime _dueDate = DateTime.now();

  @override
  void initState() {
    _dueDate = widget.dueDate;
    super.initState();
  }

  void _selectDate(BuildContext context) async {
    final DateTime? selected = await showDatePicker(
      context: context,
      initialDate: _dueDate,
      firstDate: DateTime.now(),
      lastDate: DateTime(DateTime.now().year + 100),
    );
    log.info("selected = $selected");
    if (selected != null) {
      setState(() {
        _dueDate = selected;
        widget.onChanged(_dueDate);
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(8.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Text(
            "期限",
            style: TextStyle(fontSize: 20),
          ),
          TextButton(
            child: Text(
              dateFormat.format(_dueDate),
              style: const TextStyle(fontSize: 20),
            ),
            onPressed: () => _selectDate(context),
          ),
        ],
      ),
    );
  }
}
