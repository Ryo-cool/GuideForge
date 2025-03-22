import Link from 'next/link';

export default function Home() {
  return (
    <main className="container mx-auto px-6">
      <div className="neumorph-card mt-16 text-center">
        <h1 className="mb-8 text-primary text-3xl font-bold">
          GuideForge
        </h1>
        <p className="text-xl mb-8">
          業務効率化のためのマニュアル作成・管理ツール
        </p>
        <div className="mt-12 flex justify-center gap-4">
          <Link href="/login" className="neumorph-button primary">
            ログイン
          </Link>
          <Link href="/register" className="neumorph-button">
            新規登録
          </Link>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-16">
        <div className="neumorph-card h-full">
          <h3 className="mb-4 text-xl font-semibold">簡単作成</h3>
          <p>
            直感的なインターフェースで、すばやく効率的にマニュアルを作成できます。
            ドラッグ＆ドロップでステップを並べ替えることも可能です。
          </p>
        </div>
        <div className="neumorph-card h-full">
          <h3 className="mb-4 text-xl font-semibold">画像サポート</h3>
          <p>
            各ステップに画像を添付して、視覚的に分かりやすいマニュアルを作成できます。
            手順を説明するためのスクリーンショットや図解を簡単に追加できます。
          </p>
        </div>
        <div className="neumorph-card h-full">
          <h3 className="mb-4 text-xl font-semibold">管理と共有</h3>
          <p>
            作成したマニュアルを一元管理し、必要に応じて他のユーザーと共有できます。
            カテゴリやキーワードで簡単に検索することも可能です。
          </p>
        </div>
      </div>
    </main>
  );
} 