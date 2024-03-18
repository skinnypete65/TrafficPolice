package camera

import (
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/service"
	"TrafficPolice/internal/transport/rest/dto"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"github.com/jcalabro/leb128"
	"time"
)

const (
	typeCamerus1 = "camerus1"
	typeCamerus2 = "camerus2"
	typeCamerus3 = "camerus3"

	cameraIDKey = "camera_id"
	cameraKey   = "camera"
)

type Parser struct {
	cameraService service.CameraService
}

func NewParser(cameraService service.CameraService) *Parser {
	return &Parser{
		cameraService: cameraService,
	}
}

func (p *Parser) ParseCameraInfo(payload []byte) (dto.Case, error) {
	if len(payload) == 0 {
		return dto.Case{}, errs.ErrEmptyPayload
	}

	payload = payload[2:]
	info := p.parsePayload(payload)

	var cameraType string
	var err error
	if cameraID, ok := info[cameraIDKey]; ok {
		_, err = uuid.Parse(cameraID.(string))
		if err != nil {
			return dto.Case{}, errs.ErrInvalidCameraID
		}
		cameraType, err = p.cameraService.GetCameraTypeByCameraID(cameraID.(string))
	} else if camera, ok := info[cameraKey].(map[string]any); ok {
		cameraID := camera["id"].(string)
		_, err = uuid.Parse(cameraID)
		if err != nil {
			return dto.Case{}, errs.ErrInvalidCameraID
		}
		cameraType, err = p.cameraService.GetCameraTypeByCameraID(cameraID)
	} else {
		err = errs.ErrUnknownCameraID
	}

	if err != nil {
		return dto.Case{}, err
	}

	var caseInfo dto.Case
	switch cameraType {
	case typeCamerus1:
		caseInfo, err = p.parseCamerus1(info)
	case typeCamerus2:
		caseInfo, err = p.parseCamerus2(info)
	case typeCamerus3:
		caseInfo, err = p.parseCamerus3(info)
	default:
		return dto.Case{}, errs.ErrUnknownCameraType
	}
	if err != nil {
		return dto.Case{}, err
	}

	err = p.validateCase(caseInfo)
	if err != nil {
		return dto.Case{}, err
	}
	return caseInfo, nil
}

func (p *Parser) parsePayload(payload []byte) map[string]any {
	info := make(map[string]any)

	for len(payload) > 0 {
		keySize := binary.BigEndian.Uint16(payload[:2])
		payload = payload[2:]

		valueSize := binary.BigEndian.Uint16(payload[:2])
		payload = payload[2:]

		valueType := payload[0]
		payload = payload[1:]

		keyValue := payload[:keySize]
		payload = payload[keySize:]

		value := payload[:valueSize]
		payload = payload[valueSize:]

		if valueType == 0 {
			info[string(keyValue)] = string(value)
		} else if valueType == 1 {
			n, _ := leb128.DecodeS64(bytes.NewBuffer(value))
			info[string(keyValue)] = n
		} else if valueType == 2 {
			dict := p.parsePayload(value)
			info[string(keyValue)] = dict
		}
	}

	return info
}

func (p *Parser) parseCamerus1(info map[string]any) (dto.Case, error) {
	date, err := time.Parse(time.RFC3339, info["datetime"].(string))
	if err != nil {
		return dto.Case{}, err
	}
	return dto.Case{
		Transport: dto.Transport{
			Chars:  info["transport_chars"].(string),
			Num:    info["transport_numbers"].(string),
			Region: info["transport_region"].(string),
		},
		Camera: dto.Camera{
			ID: info["camera_id"].(string),
		},
		Violation: dto.Violation{
			ID: info["violation_id"].(string),
		},
		ViolationValue: info["violation_value"].(string),
		RequiredSkill:  info["skill_value"].(int64),
		Date:           date,
	}, nil
}

func (p *Parser) parseCamerus2(info map[string]any) (dto.Case, error) {
	transport := info["transport"].(map[string]any)
	camera := info["camera"].(map[string]any)
	violation := info["violation"].(map[string]any)
	skill := info["skill"].(map[string]any)

	datetime := info["datetime"].(map[string]any)

	year := datetime["year"].(int64)
	month := datetime["month"].(int64)
	day := datetime["day"].(int64)
	hour := datetime["hour"].(int64)
	minute := datetime["minute"].(int64)
	seconds := datetime["seconds"].(int64)
	utcOffset := datetime["utc_offset"].(string)

	dateString := fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d%s", year, month, day, hour, minute, seconds, utcOffset)
	date, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		return dto.Case{}, err
	}

	return dto.Case{
		Transport: dto.Transport{
			Chars:  transport["chars"].(string),
			Num:    transport["numbers"].(string),
			Region: transport["region"].(string),
		},
		Camera: dto.Camera{
			ID: camera["id"].(string),
		},
		Violation: dto.Violation{
			ID: violation["id"].(string),
		},
		ViolationValue: violation["value"].(string),
		RequiredSkill:  skill["value"].(int64),
		Date:           date,
	}, nil

}

func (p *Parser) parseCamerus3(info map[string]any) (dto.Case, error) {
	transportStr := info["transport"].(string)
	transport := []rune(transportStr)

	chars := string(transport[1:4])
	num := string(transport[0]) + string(transport[4:6])
	region := string(transport[6:])

	camera := info["camera"].(map[string]any)
	violation := info["violation"].(map[string]any)

	return dto.Case{
		Transport: dto.Transport{
			Chars:  chars,
			Num:    num,
			Region: region,
		},
		Camera: dto.Camera{
			ID: camera["id"].(string),
		},
		Violation: dto.Violation{
			ID: violation["id"].(string),
		},
		ViolationValue: violation["value"].(string),
		RequiredSkill:  info["skill"].(int64),
		Date:           time.Unix(info["datetime"].(int64), 0),
	}, nil
}

func (p *Parser) validateCase(c dto.Case) error {
	_, err := uuid.Parse(c.Violation.ID)
	if err != nil {
		return errs.ErrInvalidViolationID
	}
	return nil
}
