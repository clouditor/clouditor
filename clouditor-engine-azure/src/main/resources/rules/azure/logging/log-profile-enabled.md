# Configure Retention Duration for Activity Log

Checks if Activity Log Retention is set for 365 days or greater.

```ccl
# TODO: retentionPolicy == 0 is also valid
Subscription has retentionPolicy.days >= 365 in all logProfiles
```

## Controls

CIS Microsoft Azure Foundations Benchmark/Azure 5.2
