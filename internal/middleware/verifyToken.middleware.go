package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/tickitz-backend/internal/jwttoken"
	"github.com/tickitz-backend/internal/repository"
	"github.com/tickitz-backend/pkg"
)

func VerifyToken(authCache *repository.AuthCacheRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) { // closure function
		token, ok := jwttoken.VerifyClientToken(ctx)
		if !ok {
			ctx.AbortWithStatus(401)
			return
		}

		var claims pkg.Claims
		if err := claims.VerifyJWT(token); err != nil {
			ctx.AbortWithStatus(401)
			return
		}

		tokenHash := jwttoken.HashToken(token)
		isActive, err := authCache.IsTokenActive(ctx.Request.Context(), tokenHash, claims.UserId)
		if !jwttoken.HandleTokenIsActive(ctx, isActive, err) {
			ctx.AbortWithStatus(401)
			return
		}

		ctx.Set("claims", claims)
		ctx.Set("token_hash", tokenHash)
		ctx.Next()
	}
}

func AuthorizeRoles(roles ...string) gin.HandlerFunc {
	allowedRoles := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowedRoles[role] = struct{}{}
	}

	return func(ctx *gin.Context) {
		claims, ok := jwttoken.GetClaims(ctx)
		if !ok {
			ctx.AbortWithStatus(401)
			return
		}

		if _, ok := allowedRoles[claims.Role]; !ok {
			ctx.AbortWithStatus(403)
			return
		}

		ctx.Next()
	}
}
