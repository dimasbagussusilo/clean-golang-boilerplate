package middlewares

// func OauthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var authModel models.JWTAccessTokenPayload

// 		cfg := config.Get()

// 		grpcAuthHost := fmt.Sprintf("%s:%d", cfg.Grpc.Client.CoreAuthService.Host, cfg.Grpc.Client.CoreAuthService.Port)

// 		var grpcConnection *grpc.ClientConn

// 		if cfg.Environment == "local" {
// 			conn, err := grpc.DialContext(c, grpcAuthHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 			if err != nil {
// 				errLog := &model.ErrorLog{
// 					SystemMessage: err.Error(),
// 					Message:       "Grpc Initial Dial failed",
// 				}

// 				response := &model.Response{
// 					Error:      errLog,
// 					StatusCode: http.StatusInternalServerError,
// 				}

// 				c.AbortWithStatusJSON(http.StatusInternalServerError, response)
// 				c.AbortWithStatusJSON(http.StatusInternalServerError, err)
// 				return
// 			}

// 			grpcConnection = conn

// 		} else {
// 			f, err := os.ReadFile(cfg.PublicSSLPath)

// 			if err != nil {
// 				errLog := &model.ErrorLog{
// 					SystemMessage: err.Error(),
// 					Message:       "Failed Read SSL file",
// 				}

// 				response := &model.Response{
// 					Error:      errLog,
// 					StatusCode: http.StatusInternalServerError,
// 				}

// 				c.AbortWithStatusJSON(http.StatusInternalServerError, response)
// 				c.AbortWithStatusJSON(http.StatusInternalServerError, err)
// 				return
// 			}

// 			p := x509.NewCertPool()
// 			p.AppendCertsFromPEM(f)
// 			tlsConfig := &tls.Config{
// 				RootCAs: p,
// 			}

// 			conn, err := grpc.DialContext(c, grpcAuthHost, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))

// 			if err != nil {
// 				errLog := &model.ErrorLog{
// 					SystemMessage: err.Error(),
// 					Message:       "Grpc Initial Dial failed",
// 				}

// 				response := &model.Response{
// 					Error:      errLog,
// 					StatusCode: http.StatusInternalServerError,
// 				}

// 				c.AbortWithStatusJSON(http.StatusInternalServerError, response)
// 				c.AbortWithStatusJSON(http.StatusInternalServerError, err)
// 				return
// 			}

// 			grpcConnection = conn
// 		}

// 		authorizationBearerToken := helper.GetAuthorizationValue(c.GetHeader("Authorization"))
// 		ctx, cancel := context.WithTimeout(c, 10*time.Second)
// 		client := pb.NewOauth2ServiceClient(grpcConnection)

// 		defer cancel()

// 		verifyAccessToken, err := client.VerifyAndValidateAccessToken(ctx, &pb.VerifyAndValidateAccessTokenRequest{Token: authorizationBearerToken})

// 		if err != nil {
// 			errLog := &model.ErrorLog{
// 				SystemMessage: err.Error(),
// 				Message:       "Grpc Dial failed Call Service",
// 			}

// 			response := &model.Response{
// 				Error:      errLog,
// 				StatusCode: http.StatusInternalServerError,
// 			}

// 			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
// 			return
// 		}

// 		if verifyAccessToken.StatusCode == 200 {
// 			authJson, _ := json.Marshal(verifyAccessToken.Data)
// 			_ = json.Unmarshal(authJson, &authModel)
// 			c.Set("user", &authModel)
// 			c.Writer.Header().Set("Content-Type", "application/json")
// 			c.Next()
// 		} else {
// 			errLog := &model.ErrorLog{
// 				SystemMessage: verifyAccessToken.Error.SystemMessage,
// 				Message:       verifyAccessToken.Error.Message,
// 			}

// 			response := &model.Response{
// 				Error:      errLog,
// 				StatusCode: int(verifyAccessToken.StatusCode),
// 			}

// 			c.AbortWithStatusJSON(int(verifyAccessToken.StatusCode), response)
// 			return
// 		}
// 	}
// }
