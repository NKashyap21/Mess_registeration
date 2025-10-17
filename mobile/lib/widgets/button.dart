import 'package:flutter/material.dart';
import 'package:flutter_frontend/theme/app_theme.dart';

class AppButton extends StatelessWidget {
  final String label;
  final VoidCallback onPressed;
  final Widget? icon;
  final Color? backgroundColor;
  final Color? textColor;
  final EdgeInsets? padding;
  final bool isOutlined;
  final Color? borderColor;
  final double? width;
  final BorderRadius borderRadius;

  const AppButton({
    Key? key,
    required this.label,
    required this.onPressed,
    this.icon,
    this.backgroundColor,
    this.textColor,
    this.padding,
    this.isOutlined = false,
    this.borderColor,
    this.width,
    this.borderRadius = const BorderRadius.all(Radius.circular(0)),
  }) : super(key: key);

  factory AppButton.icon({
    required String label,
    required VoidCallback onPressed,
    required Widget icon,
    Color? backgroundColor,
    Color? textColor,
    EdgeInsets? padding,
    bool isOutlined = false,
    Color? borderColor,
    double? width,
    required BorderRadius borderRadius,
  }) {
    return AppButton(
      label: label,
      onPressed: onPressed,
      icon: icon,
      backgroundColor: backgroundColor,
      textColor: textColor,
      padding: padding,
      isOutlined: isOutlined,
      borderColor: borderColor,
      width: width,
      borderRadius: borderRadius,
    );
  }

  @override
  Widget build(BuildContext context) {
    final ButtonStyle buttonStyle = isOutlined
        ? OutlinedButton.styleFrom(
            side: BorderSide(color: borderColor ?? AppColors.midGrey),
            padding:
                padding ??
                const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            backgroundColor: backgroundColor,
            foregroundColor: textColor,
            minimumSize: width != null ? Size(width!, 0) : null,
            shape: RoundedRectangleBorder(borderRadius: borderRadius),
          )
        : ElevatedButton.styleFrom(
            backgroundColor: backgroundColor,
            foregroundColor: textColor,
            padding:
                padding ??
                const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            minimumSize: width != null ? Size(width!, 0) : null,
            shape: RoundedRectangleBorder(borderRadius: borderRadius),
          );

    final Widget buttonChild = Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        if (icon != null) ...[
          Container(margin: const EdgeInsets.only(right: 8.0), child: icon),
        ],
        Text(label),
      ],
    );

    return isOutlined
        ? OutlinedButton(
            style: buttonStyle,
            onPressed: onPressed,
            child: buttonChild,
          )
        : ElevatedButton(
            style: buttonStyle,
            onPressed: onPressed,
            child: buttonChild,
          );
  }
}
