import 'package:flutter/material.dart';
import 'package:flutter_frontend/theme/app_theme.dart';
import 'package:flutter_frontend/widgets/button.dart';
import 'package:flutter_frontend/screens/registration.dart';
import 'package:flutter_frontend/widgets/profile_button.dart';
import 'package:flutter_frontend/screens/swap.dart';
import 'package:flutter_frontend/api/api_service.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  Map<String, dynamic>? userInfo;
  Map<String, dynamic>? messInfo;
  Map<String, dynamic>? registrationStatus;
  bool isLoading = true;
  bool isLoadingMess = true;
  bool isLoadingRegistrationStatus = true;
  bool isRegisteringVeg = false;
  String? error;

  @override
  void initState() {
    super.initState();
    _loadUserInfo();
    _loadMessInfo();
    _loadRegistrationStatus();
  }

  Future<void> _loadUserInfo() async {
    try {
      final response = await ApiService.getUserInfo();
      setState(() {
        if (response['error'] != null) {
          error = response['error'];
        } else {
          userInfo = response['data'];
        }
        isLoading = false;
      });
    } catch (e) {
      setState(() {
        error = 'Failed to load user info: $e';
        isLoading = false;
      });
    }
  }

  Future<void> _loadMessInfo() async {
    try {
      final response = await ApiService.getMess();
      setState(() {
        if (response['error'] != null) {
          messInfo = null;
        } else {
          messInfo = response;
        }
        isLoadingMess = false;
      });
    } catch (e) {
      setState(() {
        messInfo = null;
        isLoadingMess = false;
      });
    }
  }

  Future<void> _loadRegistrationStatus() async {
    try {
      final response = await ApiService.getRegistrationStatus();

      print(response);

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
    }
  }

  Future<void> _refreshData() async {
    setState(() {
      isLoading = true;
      isLoadingMess = true;
      isLoadingRegistrationStatus = true;
      error = null;
    });

    await Future.wait([
      _loadUserInfo(),
      _loadMessInfo(),
      _loadRegistrationStatus(),
    ]);
  }

  void _showVegMessConfirmationDialog() {
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
              const Text(
                'You are Registering for\nVegetarian Mess',
                textAlign: TextAlign.center,
                style: TextStyle(
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
                  onPressed: () => _registerForVegMess(dialogContext),
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

  Future<void> _registerForVegMess(BuildContext dialogContext) async {
    Navigator.of(dialogContext).pop();

    setState(() {
      isRegisteringVeg = true;
    });

    try {
      final response = await ApiService.registerVegMess();

      setState(() {
        isRegisteringVeg = false;
      });

      if (response['error'] != null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Veg Mess registration failed: ${response['error']}'),
            backgroundColor: AppColors.orange,
            duration: const Duration(seconds: 3),
          ),
        );
      } else {
        _loadUserInfo();
        _loadMessInfo();
        _loadRegistrationStatus();

        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Successfully registered for Vegetarian Mess'),
            backgroundColor: Colors.green,
            duration: Duration(seconds: 2),
          ),
        );
      }
    } catch (e) {
      setState(() {
        isRegisteringVeg = false;
      });

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Veg Mess registration failed: $e'),
          backgroundColor: AppColors.orange,
          duration: const Duration(seconds: 3),
        ),
      );
    }
  }

  bool get isRegularRegistrationOpen {
    return registrationStatus?['regular'] ?? false;
  }

  bool get isVegRegistrationOpen {
    return registrationStatus?['veg'] ?? false;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: LayoutBuilder(
          builder: (context, constraints) {
            return RefreshIndicator(
              onRefresh: _refreshData,
              child: SingleChildScrollView(
                physics: const AlwaysScrollableScrollPhysics(),
                child: ConstrainedBox(
                  constraints: BoxConstraints(minHeight: constraints.maxHeight),
                  child: IntrinsicHeight(
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 16.0),
                      child: Column(
                        children: [
                          const SizedBox(height: 20),

                          Align(
                            alignment: Alignment.centerRight,
                            child: ProfileButton(userInfo: userInfo),
                          ),

                          const Spacer(),

                          (isLoading ||
                                  isLoadingMess ||
                                  isLoadingRegistrationStatus)
                              ? const Center(child: CircularProgressIndicator())
                              : error != null
                              ? Center(
                                  child: Column(
                                    children: [
                                      Text(
                                        'Error: $error',
                                        style: TextStyle(
                                          color: AppColors.orange,
                                        ),
                                      ),
                                      const SizedBox(height: 10),
                                      ElevatedButton(
                                        onPressed: _loadUserInfo,
                                        child: const Text('Retry'),
                                      ),
                                    ],
                                  ),
                                )
                              : _buildMainCard(context),

                          const Spacer(),

                          _buildFooter(),

                          const SizedBox(height: 16),
                        ],
                      ),
                    ),
                  ),
                ),
              ),
            );
          },
        ),
      ),
    );
  }

  Widget _buildMainCard(BuildContext context) {
    String messDisplayName;
    bool isRegistered = false;

    if (messInfo != null) {
      final messId = messInfo!['mess'];
      final status = messInfo!['status'];

      if (messId != null && messId != 0) {
        isRegistered = true;

        if (status == 'pending_sync') {
          String messName;

          switch (messId) {
            case 1:
              messName = "Mess A LDH";
              break;
            case 2:
              messName = "Mess A UDH";
              break;
            case 3:
              messName = "Mess B LDH";
              break;
            case 4:
              messName = "Mess B UDH";
              break;
            case 5:
              messName = "Veg Mess";
            default:
              messName = "Mess $messId";
          }

          messDisplayName = messName;
        } else {
          messDisplayName = messInfo!['mess_name'] ?? 'Mess ${messId}';
        }
      } else {
        messDisplayName = 'Unregistered';
      }
    } else {
      messDisplayName = 'Unregistered';
    }

    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(color: AppColors.darkGrey),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          // Icon and title
          Image.asset(
            'assets/images/mess_icon.png',
            width: 40,
            height: 40,
            fit: BoxFit.contain,
          ),
          const SizedBox(height: 8),
          const Text(
            'Mess Portal',
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.bold,
              color: AppColors.white,
            ),
          ),
          const SizedBox(height: 32),

          Container(
            width: double.infinity,
            child: _buildInfoRow(
              'Current Registered Mess :',
              messDisplayName,
              !isRegistered,
            ),
          ),
          const SizedBox(height: 32),

          if (isVegRegistrationOpen) ...[
            AppButton(
              label: isRegisteringVeg
                  ? 'Registering...'
                  : 'Register for Veg Mess',
              onPressed: isRegisteringVeg
                  ? () {}
                  : _showVegMessConfirmationDialog,
              backgroundColor: Colors.transparent,
              textColor: isRegisteringVeg
                  ? AppColors.lightGrey
                  : AppColors.white,
              borderColor: isRegisteringVeg
                  ? AppColors.lightGrey
                  : AppColors.orange,
              isOutlined: true,
              width: double.infinity,
              borderRadius: BorderRadius.circular(0),
            ),
            const SizedBox(height: 10),
          ] else ...[
            Container(
              width: double.infinity,
              padding: const EdgeInsets.symmetric(vertical: 12),
              decoration: BoxDecoration(
                border: Border.all(color: AppColors.lightGrey),
                borderRadius: BorderRadius.circular(0),
              ),
              child: const Text(
                'Veg Registration Closed',
                textAlign: TextAlign.center,
                style: TextStyle(color: AppColors.lightGrey, fontSize: 14),
              ),
            ),
            const SizedBox(height: 10),
          ],

          if (isRegularRegistrationOpen) ...[
            AppButton(
              label: 'Go for Registration',
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (context) => const RegistrationScreen(),
                  ),
                ).then((_) {
                  _loadUserInfo();
                  _loadMessInfo();
                });
              },
              backgroundColor: Colors.transparent,
              textColor: AppColors.white,
              borderColor: AppColors.midGrey,
              isOutlined: true,
              width: double.infinity,
              borderRadius: BorderRadius.circular(0),
            ),
            const SizedBox(height: 10),
          ] else ...[
            Container(
              width: double.infinity,
              padding: const EdgeInsets.symmetric(vertical: 12),
              decoration: BoxDecoration(
                border: Border.all(color: AppColors.lightGrey),
                borderRadius: BorderRadius.circular(0),
              ),
              child: const Text(
                'Regular Registration Closed',
                textAlign: TextAlign.center,
                style: TextStyle(color: AppColors.lightGrey, fontSize: 14),
              ),
            ),
            const SizedBox(height: 10),
          ],

          AppButton(
            label: 'Mess Swap',
            onPressed: () {
              if (isRegistered) {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => const SwapScreen()),
                );
              } else {
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(
                    content: Text('Please register for a mess first'),
                    backgroundColor: AppColors.orange,
                  ),
                );
              }
            },
            backgroundColor: Colors.transparent,
            textColor: isRegistered ? AppColors.white : AppColors.lightGrey,
            borderColor: isRegistered ? AppColors.midGrey : AppColors.lightGrey,
            isOutlined: true,
            width: double.infinity,
            borderRadius: BorderRadius.circular(0),
          ),
        ],
      ),
    );
  }

  Widget _buildInfoRow(String label, String value, bool isError) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Flexible(
          flex: 3,
          child: Text(
            label,
            style: const TextStyle(color: AppColors.white, fontSize: 14),
          ),
        ),
        const SizedBox(width: 8),
        Flexible(
          flex: 2,
          child: Text(
            value,
            style: TextStyle(
              color: isError ? AppColors.orange : AppColors.white,
              fontSize: 14,
              fontWeight: isError ? FontWeight.bold : FontWeight.normal,
            ),
            textAlign: TextAlign.right,
            overflow: TextOverflow.ellipsis,
          ),
        ),
      ],
    );
  }

  Widget _buildFooter() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: const [
        Text(
          'Brought to you with ',
          style: TextStyle(color: AppColors.lightGrey),
        ),
        Icon(Icons.favorite, color: Colors.deepPurple, size: 16),
        Text(' by Lambda', style: TextStyle(color: AppColors.lightGrey)),
      ],
    );
  }
}
