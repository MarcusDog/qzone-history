// export_usecase.go

package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/domain/repository"
	"qzone-history/internal/domain/usecase"
	"regexp"
)

var qqNumberPattern = regexp.MustCompile(`^[0-9]+$`)

type exportUseCase struct {
	momentRepo       repository.MomentRepository
	boardMessageRepo repository.BoardMessageRepository
	friendRepo       repository.FriendRepository
}

func NewExportUseCase(
	momentRepo repository.MomentRepository,
	boardMessageRepo repository.BoardMessageRepository,
	friendRepo repository.FriendRepository,
) usecase.ExportUseCase {
	return &exportUseCase{
		momentRepo:       momentRepo,
		boardMessageRepo: boardMessageRepo,
		friendRepo:       friendRepo,
	}
}

func (u *exportUseCase) ExportUserDataToJSON(ctx context.Context, userQQ string) error {
	if !qqNumberPattern.MatchString(userQQ) {
		return fmt.Errorf("invalid user QQ: expected digits only")
	}

	// 获取用户数据。任何一项失败都应中止，避免生成看似成功但不完整的备份。
	moments, err := u.momentRepo.FindByUserQQ(ctx, userQQ, -1, 0)
	if err != nil {
		return fmt.Errorf("failed to load moments: %w", err)
	}
	boardMessages, err := u.boardMessageRepo.FindByUserQQ(ctx, userQQ, -1, 0)
	if err != nil {
		return fmt.Errorf("failed to load board messages: %w", err)
	}
	friends, err := u.friendRepo.FindFriendsByUserQQ(ctx, userQQ)
	if err != nil {
		return fmt.Errorf("failed to load friends: %w", err)
	}

	// 创建导出数据结构
	exportData := struct {
		Moments       []entity.Moment       `json:"moments"`
		BoardMessages []entity.BoardMessage `json:"boardMessages"`
		Friends       []entity.Friend       `json:"friends"`
	}{
		Moments:       moments,
		BoardMessages: boardMessages,
		Friends:       friends,
	}

	// 转换为JSON
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// 写入文件
	filename := fmt.Sprintf("%s_export.json", userQQ)
	err = os.WriteFile(filename, jsonData, 0600)
	if err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

func (u *exportUseCase) ExportUserDataToExcel(ctx context.Context, userQQ string) error {
	// TODO 实现Excel导出逻辑
	return fmt.Errorf("ExportUserDataToExcel not implemented")
}

func (u *exportUseCase) ExportUserDataToHTML(ctx context.Context, userQQ string) error {
	// TODO 实现HTML导出逻辑
	return fmt.Errorf("ExportUserDataToHTML not implemented")
}
