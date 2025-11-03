package fixtures

// TestConfig contains test configuration data.
const TestConfig = `
current-context: test
contexts:
  test:
    brokers:
      - localhost:9092
    schema-registry: http://localhost:8081
  prod:
    brokers:
      - prod-kafka-1:9092
      - prod-kafka-2:9092
    schema-registry: https://schema-registry.prod.example.com
`

// TestTopicsList contains mock kafkactl output for topics.
const TestTopicsList = `
- name: test-topic-1
  partitions: 3
  replication-factor: 1
- name: test-topic-2
  partitions: 6
  replication-factor: 2
`

// TestSchemasList contains mock kafkactl output for schemas.
const TestSchemasList = `
- subject: test-subject-value
  version: 1
  schema: '{"type":"record","name":"Test","fields":[]}'
`
