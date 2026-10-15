#ifndef SPACESEE_MOCK_SHIM_H
#define SPACESEE_MOCK_SHIM_H

#ifdef __cplusplus
extern "C" {
#endif

int SPACESEE_mock_load(const char *workDir, const char *assetsDir);
const char *SPACESEE_mock_loadError(void);
void SPACESEE_mock_unload(void);
int SPACESEE_mock_isHeavy(void);
int SPACESEE_mock_assetCount(void);
const char *SPACESEE_mock_asset(int index);
const char *SPACESEE_mock_version(void);
const char *SPACESEE_mock_extension(void);
const char *SPACESEE_mock_filetypes(void);
long long SPACESEE_mock_maxInputSize(void);
long long SPACESEE_mock_analyze(const char *inputPath, const char *outputPath, const char *dataDir);
int SPACESEE_mock_progress(long long jobId, int scale);
const char *SPACESEE_mock_progressString(long long jobId);
void SPACESEE_mock_free(void *ptr);

#ifdef __cplusplus
}
#endif

#endif
