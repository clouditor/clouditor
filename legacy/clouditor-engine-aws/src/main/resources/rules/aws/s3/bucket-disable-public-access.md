# Block Public Access For Buckets

By default, public access for buckets should be blocked.

```ccl
Bucket has publicAccessBlockConfiguration.blockPublicAcls == true
Bucket has publicAccessBlockConfiguration.blockPublicPolicy == true
Bucket has publicAccessBlockConfiguration.ignorePublicAcls == true
Bucket has publicAccessBlockConfiguration.restrictPublicBuckets == true
```
