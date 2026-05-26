# REVIEW: real ntfy ping

Goal: Confirm a real ntfy.sh topic receives a test notification.

Context:
- Verification task; needs network.

Do:
1. With a chosen topic, call Publish once (scratch harness) and confirm receipt on phone/web.
2. Confirm a wrong token/topic fails gracefully without crashing.

Done when:
- Real notification received.
- Failure path graceful.
- Tests pass.
