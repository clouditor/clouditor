# Rotate Access Keys Regularly admin

Checks if AWS access keys are rotated regularly and are not older than a specified amount of time, i.e. 1000 days.

```ccl
User with userName=="admin" has createDate younger 90 days in all accessKeys
```

## Controls

* "IAM-02"
* "AWS 1.4"
* "BSI C5/KRY-04"
