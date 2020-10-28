package storage

import (
	"database/sql"
	"encoding/json"
	"strings"

	"git.gdpteam.com/4gen/4g-tags-api/db/rds"
	"git.gdpteam.com/4gen/4g-tags-api/pkg/models/tag"
	"github.com/go-redis/redis/v8"
)

const redisVar = "tags"

// RedisTag used for work with Redis - tag
type RedisTag struct {
	db  *sql.DB
	rdb *rds.Rds
}

// NewRedisTag return a new pointer of RedisTag
func NewRedisTag(db *sql.DB, rdb *rds.Rds) *RedisTag {
	return &RedisTag{db, rdb}
}

// Create implement the interface tag.Storage for Redis
func (t *RedisTag) Create(m *tag.Model) error {
	storageTag := NewMySQLTag(t.db, nil)
	serviceTag := tag.NewService(storageTag)
	if err := serviceTag.Create(m); err != nil {
		return err
	}

	if err := t.RefreshData(); err != nil {
		return err
	}

	return nil
}

// GetAll implement the interface tag.Storage for Redis
func (t *RedisTag) GetAll() (tag.Models, error) {
	models := make(tag.Models, 0)

	redisValue, err := t.getRedisTagsMap()
	if err != nil {
		return nil, err
	}

	for _, tagVal := range redisValue {
		tagVal.Name = strings.Title(strings.ToLower(tagVal.Name))
		models = append(models, tagVal)
	}

	return models, nil
}

// GetAllWithFields implement the interface tag.Storage for Redis
func (t *RedisTag) GetAllWithFields(fields ...string) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)

	redisValue, err := t.getRedisTagsMap()
	if err != nil {
		return nil, err
	}

	valuesMap := make(map[string]interface{})

	for _, tagVal := range redisValue {
		value := make(map[string]interface{})

		valuesMap["tagId"] = tagVal.TagID
		valuesMap["name"] = tagVal.Name
		valuesMap["insUserId"] = tagVal.InsUserID
		valuesMap["insDate"] = tagVal.InsDate
		valuesMap["insDatetime"] = tagVal.InsDateTime
		valuesMap["insTimestamp"] = tagVal.InsTimestamp

		for _, field := range fields {
			if val, ok := valuesMap[field]; ok {
				value[field] = val
			}
		}

		values = append(values, value)
	}

	return values, nil
}

// GetByID implement the interface tag.Storage for Redis
func (t *RedisTag) GetByID(tagID uint) (*tag.Model, error) {
	redisValue, err := t.getRedisTagsMap()
	if err != nil {
		return nil, err
	}

	tagValue, ok := redisValue[tagID]
	if !ok {
		return nil, redis.Nil
	}

	tagValue.Name = strings.Title(strings.ToLower(tagValue.Name))

	return tagValue, nil
}

// GetByIDs implement the interface tag.Storage
func (t *RedisTag) GetByIDs(tagIDs ...uint) (tag.Models, error) {
	models := make(tag.Models, 0)
	redisValue, err := t.getRedisTagsMap()
	if err != nil {
		return nil, err
	}

	for _, tagID := range tagIDs {
		tagVal, ok := redisValue[tagID]
		if ok {
			models = append(models, tagVal)
		}
	}

	return models, nil
}

// Update implements the interface tag.Storage for Redis
func (t *RedisTag) Update(m *tag.Model) error {
	storageAudit := NewMySQLTagAudit(t.db)
	storageTag := NewMySQLTag(t.db, storageAudit)

	servParam := tag.NewService(storageTag)
	if err := servParam.Update(m); err != nil {
		return err
	}

	if err := t.RefreshData(); err != nil {
		return err
	}

	return nil
}

// RefreshData refreshes the redis storage
func (t *RedisTag) RefreshData() error {
	storageTag := NewMySQLTag(t.db, nil)
	serviceTag := tag.NewService(storageTag)
	ms, err := serviceTag.GetAll()
	if err != nil {
		return err
	}

	redisValue := make(map[uint]*tag.Model)

	for _, m := range ms {
		redisValue[m.TagID] = m
	}

	j, err := json.Marshal(redisValue)
	if err != nil {
		return err
	}

	err = t.rdb.InfoToRAM(redisVar, string(j))

	return nil
}

func (t *RedisTag) getRedisTagsMap() (map[uint]*tag.Model, error) {
	redisValue := make(map[uint]*tag.Model)

	tagJSONstring, err := t.rdb.InfoFromRAM(redisVar)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(tagJSONstring), &redisValue); err != nil {
		return nil, err
	}

	return redisValue, nil
}

// CleanData cleans redis storage
func (t *RedisTag) CleanData() error {
	return t.rdb.DelFromRAM(redisVar)
}
