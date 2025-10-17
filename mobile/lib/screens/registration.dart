import 'package:flutter/material.dart';
import 'package:flutter_frontend/theme/app_theme.dart';
import 'package:flutter_frontend/widgets/button.dart';
import 'package:flutter_frontend/screens/home.dart';
import 'package:flutter_frontend/widgets/profile_button.dart';
import 'package:flutter_frontend/api/api_service.dart';

class RegistrationScreen extends StatefulWidget {
  const RegistrationScreen({super.key});

  @override
  State<RegistrationScreen> createState() => _RegistrationScreenState();
}

class _RegistrationScreenState extends State<RegistrationScreen> {
  String? selectedMess;
  bool isRegistering = false;
  bool isLoadingStats = true;
  bool isLoadingRegistrationStatus = true;
  Map<String, dynamic>? registrationStatus;

  Map<String, Map<String, dynamic>> messStats = {};

  int _getMessNumber(String messKey) {
    switch (messKey) {
      case 'Mess-A_LDH':
        return 1;
      case 'Mess-A_UDH':
        return 2;
      case 'Mess-B_LDH':
        return 3;
      case 'Mess-B_UDH':
        return 4;
      default:
        return 1;
    }
  }

  @override
  void initState() {
    super.initState();
    _loadMessStats();
    _loadRegistrationStatus();
  }

  Future<void> _loadMessStats() async {
    try {
      setState(() {
        isLoadingStats = true;
      });

      final response = await ApiService.getMessStats();

      if (response['error'] == null && response['data'] != null) {
        final data = response['data'] as Map<String, dynamic>;

        print(data);

        setState(() {
          messStats = {
            'Mess-A_LDH': {
              'current': data['mess1  ']?['registered'] ?? 0,
              'capacity': data['mess1']?['capacity'] ?? 0,
            },
            'Mess-A_UDH': {
              'current': data['mess2']?['registered'] ?? 0,
              'capacity': data['mess2']?['capacity'] ?? 0,
            },
            'Mess-B_LDH': {
              'current': data['mess3']?['registered'] ?? 0,
              'capacity': data['mess3']?['capacity'] ?? 0,
            },
            'Mess-B_UDH': {
              'current': data['mess4']?['registered'] ?? 0,
              'capacity': data['mess4']?['capacity'] ?? 0,
            },
          };
          isLoadingStats = false;
        });
      } else {
        setState(() {
          messStats = {
            'Mess-A_LDH': {'current': 0, 'capacity': 2500},
            'Mess-A_UDH': {'current': 0, 'capacity': 2500},
            'Mess-B_LDH': {'current': 0, 'capacity': 2500},
            'Mess-B_UDH': {'current': 0, 'capacity': 2500},
          };
          isLoadingStats = false;
        });
      }
    } catch (e) {
      setState(() {
        messStats = {
          'Mess-A_LDH': {'current': 0, 'capacity': 2500},
          'Mess-A_UDH': {'current': 0, 'capacity': 2500},
          'Mess-B_LDH': {'current': 0, 'capacity': 2500},
          'Mess-B_UDH': {'current': 0, 'capacity': 2500},
        };
        isLoadingStats = false;
      });
      print('Error loading mess stats: $e');
    }
  }

  Future<void> _loadRegistrationStatus() async {
    try {
      setState(() {
        isLoadingRegistrationStatus = true;
      });

      final response = await ApiService.getRegistrationStatus();

      setState(() {
        if (response['error'] != null) {
          registrationStatus = null;
        } else {
          registrationStatus = response;
        }
        isLoadingRegistrationStatus = false;
      });
    } catch (e) {
      setState(() {
        registrationStatus = null;
        isLoadingRegistrationStatus = false;
      });
      print('Error loading registration status: $e');
    }
  }

  bool get isRegularRegistrationOpen {
    return registrationStatus?['regular'] ?? false;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        title: const Text(
          'Mess Registration',
          style: TextStyle(
            color: AppColors.white,
            fontSize: 22,
            fontWeight: FontWeight.normal,
          ),
        ),
        centerTitle: true,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: AppColors.white),
          onPressed: () => Navigator.pop(context),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh, color: AppColors.white),
            onPressed: (isLoadingStats || isLoadingRegistrationStatus)
                ? null
                : () {
                    _loadMessStats();
                    _loadRegistrationStatus();
                  },
          ),
          const Padding(
            padding: EdgeInsets.only(right: 16.0),
            child: ProfileButton(userInfo: null),
          ),
        ],
      ),
      body: Column(
        children: [
          Expanded(
            child: Container(
              margin: const EdgeInsets.all(16.0),
              decoration: BoxDecoration(
                color: AppColors.darkGrey,
                borderRadius: BorderRadius.circular(12.0),
              ),
              padding: const EdgeInsets.symmetric(
                horizontal: 24.0,
                vertical: 32.0,
              ),
              child: (isLoadingStats || isLoadingRegistrationStatus)
                  ? const Center(
                      child: CircularProgressIndicator(
                        valueColor: AlwaysStoppedAnimation<Color>(
                          AppColors.orange,
                        ),
                      ),
                    )
                  : !isRegularRegistrationOpen
                  ? Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          const Icon(
                            Icons.schedule,
                            size: 64,
                            color: AppColors.lightGrey,
                          ),
                          const SizedBox(height: 16),
                          const Text(
                            'Registration Closed',
                            style: TextStyle(
                              color: AppColors.white,
                              fontSize: 24,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 8),
                          const Text(
                            'Mess registration is currently not available.\nPlease check back during the registration period.',
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              color: AppColors.lightGrey,
                              fontSize: 16,
                            ),
                          ),
                        ],
                      ),
                    )
                  : SingleChildScrollView(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          _buildMessSection('Mess-A', ['LDH', 'UDH']),
                          const SizedBox(height: 60),
                          _buildMessSection('Mess-B', ['LDH', 'UDH']),
                        ],
                      ),
                    ),
            ),
          ),
          Container(
            width: 200,
            height: 45,
            margin: const EdgeInsets.only(bottom: 32.0),
            child: AppButton(
              label: isRegistering ? 'Registering...' : 'Register',
              onPressed: (isRegistering || !isRegularRegistrationOpen)
                  ? () {}
                  : () => _showConfirmationDialog(),
              backgroundColor: (!isRegularRegistrationOpen)
                  ? AppColors.lightGrey
                  : AppColors.darkGrey,
              textColor: AppColors.white,
              width: double.infinity,
              padding: EdgeInsets.zero,
              borderRadius: BorderRadius.circular(8.0),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMessSection(String messName, List<String> divisions) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          messName,
          style: const TextStyle(
            color: AppColors.white,
            fontSize: 24,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 16),
        ...divisions.map((division) {
          final optionKey = '$messName\_$division';
          final current = messStats[optionKey]?['current'] ?? 0;
          final capacity = messStats[optionKey]?['capacity'] ?? 0;
          final progressPercentage = capacity > 0 ? (current / capacity) : 0.0;

          return Padding(
            padding: const EdgeInsets.only(bottom: 30.0),
            child: InkWell(
              onTap: () {
                setState(() {
                  selectedMess = optionKey;
                });
              },
              splashColor: Colors.transparent,
              highlightColor: Colors.transparent,
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  SizedBox(
                    width: 60,
                    child: Text(
                      division,
                      style: const TextStyle(
                        color: AppColors.white,
                        fontSize: 16,
                      ),
                    ),
                  ),
                  Expanded(
                    child: Container(
                      height: 10,
                      decoration: BoxDecoration(
                        color: AppColors.midGrey,
                        borderRadius: BorderRadius.circular(5),
                      ),
                      child: FractionallySizedBox(
                        alignment: Alignment.centerLeft,
                        widthFactor: progressPercentage.clamp(0.0, 1.0),
                        child: Container(
                          decoration: BoxDecoration(
                            color: progressPercentage < 0.8
                                ? AppColors.orange
                                : Colors.red,
                            borderRadius: BorderRadius.circular(5),
                          ),
                        ),
                      ),
                    ),
                  ),
                  Container(
                    width: 100,
                    padding: const EdgeInsets.symmetric(horizontal: 12.0),
                    child: Text(
                      '$current/$capacity',
                      style: const TextStyle(
                        color: AppColors.white,
                        fontSize: 14,
                      ),
                      textAlign: TextAlign.center,
                    ),
                  ),
                  Container(
                    width: 24,
                    height: 24,
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      border: Border.all(
                        color: selectedMess == optionKey
                            ? AppColors.orange
                            : AppColors.white,
                        width: 1.5,
                      ),
                      color: selectedMess == optionKey
                          ? AppColors.white
                          : Colors.transparent,
                    ),
                    child: selectedMess == optionKey
                        ? const Center(
                            child: Icon(
                              Icons.circle,
                              size: 12,
                              color: AppColors.black,
                            ),
                          )
                        : null,
                  ),
                ],
              ),
            ),
          );
        }),
      ],
    );
  }

  void _showConfirmationDialog() {
    if (selectedMess == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Please select a mess first'),
          backgroundColor: AppColors.orange,
        ),
      );
      return;
    }

    final messNameParts = selectedMess!.split('_');
    final messDisplayName = '${messNameParts[0]} - ${messNameParts[1]}';

    showDialog(
      context: context,
      barrierDismissible: true,
      builder: (dialogContext) => Dialog(
        backgroundColor: Colors.black,
        insetPadding: const EdgeInsets.symmetric(horizontal: 40),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8.0)),
        child: Padding(
          padding: const EdgeInsets.all(20.0),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Align(
                alignment: Alignment.topRight,
                child: InkWell(
                  onTap: () => Navigator.of(dialogContext).pop(),
                  child: const Icon(Icons.close, color: AppColors.white),
                ),
              ),
              const SizedBox(height: 16),
              Text(
                'You are Registering for\n$messDisplayName',
                textAlign: TextAlign.center,
                style: const TextStyle(
                  color: AppColors.white,
                  fontSize: 16,
                  height: 1.5,
                ),
              ),
              const SizedBox(height: 24),
              SizedBox(
                width: 120,
                height: 40,
                child: AppButton(
                  label: 'Confirm',
                  onPressed: () =>
                      _registerForMess(dialogContext, messDisplayName),
                  backgroundColor: AppColors.darkGrey,
                  textColor: AppColors.white,
                  borderRadius: BorderRadius.circular(8.0),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Future<void> _registerForMess(
    BuildContext dialogContext,
    String messDisplayName,
  ) async {
    Navigator.of(dialogContext).pop();

    setState(() {
      isRegistering = true;
    });

    try {
      final messNumber = _getMessNumber(selectedMess!);
      final response = await ApiService.registerMess(mess: messNumber);

      setState(() {
        isRegistering = false;
      });

      if (response['error'] != null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Registration failed: ${response['error']}'),
            backgroundColor: AppColors.orange,
            duration: const Duration(seconds: 3),
          ),
        );
      } else {
        _loadMessStats();

        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Successfully registered for $messDisplayName'),
            backgroundColor: Colors.green,
            duration: const Duration(seconds: 2),
          ),
        );

        Future.delayed(const Duration(milliseconds: 500), () {
          Navigator.pushAndRemoveUntil(
            context,
            MaterialPageRoute(builder: (context) => const HomeScreen()),
            (route) => false,
          );
        });
      }
    } catch (e) {
      setState(() {
        isRegistering = false;
      });

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Registration failed: $e'),
          backgroundColor: AppColors.orange,
          duration: const Duration(seconds: 3),
        ),
      );
    }
  }
}
