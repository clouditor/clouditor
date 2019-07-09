# Configure Password Policy According to Best-Practices

Checks the AWS account password policy.

```ccl
PasswordPolicy has requireUppercaseCharacters == true
PasswordPolicy has requireLowercaseCharacters == true
PasswordPolicy has requireNumbers == true
PasswordPolicy has minimumPasswordLength >= 14
PasswordPolicy has passwordReusePrevention >= 24
PasswordPolicy has maxPasswordAge <= 90
```

## Controls
  
* AWS 1.5
* AWS 1.6
* AWS 1.7
* AWS 1.8
* AWS 1.9
* AWS 1.10
* AWS 1.11
* BSI C5/IDM-11
